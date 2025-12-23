import 'package:flutter/material.dart';

import '../../api/api.dart';
import '../../models/task.dart';
import '../../widgets/app_loading_indicator.dart';
import '../../widgets/app_snack_bar.dart';

import '../reminders/reminder_preset_sheet.dart';

import 'tasks_api.dart';

class TasksScreen extends StatefulWidget {
  const TasksScreen({super.key});

  @override
  State<TasksScreen> createState() => _TasksScreenState();
}

class _TasksScreenState extends State<TasksScreen> {
  late Future<List<Task>> _future;

  final _titleController = TextEditingController();
  final _titleFocusNode = FocusNode();

  bool _isCreating = false;
  final _updatingTaskIds = <String>{};

  @override
  void initState() {
    super.initState();
    _future = TasksApi().fetchTasks();
  }

  Future<void> _toggleDone(Task task) async {
    if (task.id.isEmpty) {
      showAppErrorSnackBar(context, 'タスクIDが不正です');
      return;
    }
    if (_updatingTaskIds.contains(task.id)) return;

    setState(() {
      _updatingTaskIds.add(task.id);
    });

    try {
      await TasksApi().updateTask(id: task.id, isDone: !task.isDone);
      _reload();
    } on ApiException catch (e) {
      if (mounted) {
        showAppErrorSnackBar(context, e.message);
      }
    } catch (_) {
      if (mounted) {
        showAppErrorSnackBar(context, 'タスクの更新に失敗しました');
      }
    } finally {
      if (mounted) {
        setState(() {
          _updatingTaskIds.remove(task.id);
        });
      }
    }
  }

  String? _formatDateTime(DateTime? dt) {
    if (dt == null) return null;
    final y = dt.year.toString().padLeft(4, '0');
    final m = dt.month.toString().padLeft(2, '0');
    final d = dt.day.toString().padLeft(2, '0');
    final hh = dt.hour.toString().padLeft(2, '0');
    final mm = dt.minute.toString().padLeft(2, '0');
    return '$y-$m-$d $hh:$mm';
  }

  Future<void> _editTask(Task task) async {
    if (task.id.isEmpty) {
      showAppErrorSnackBar(context, 'タスクIDが不正です');
      return;
    }
    if (_updatingTaskIds.contains(task.id)) return;

    final controller = TextEditingController(text: task.title);
    DateTime? dueAt = task.dueAt;
    DateTime? scheduledAt = task.scheduledAt;

    final result = await showDialog<_TaskEditResult>(
      context: context,
      builder: (context) {
        return StatefulBuilder(
          builder: (context, setState) {
            return AlertDialog(
              title: const Text('タスク編集'),
              content: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  TextField(
                    controller: controller,
                    decoration: const InputDecoration(labelText: 'タイトル'),
                    textInputAction: TextInputAction.done,
                  ),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        child: Text(_formatDateTime(dueAt) ?? '期限なし'),
                      ),
                      TextButton(
                        onPressed: () async {
                          final pickedDate = await showDatePicker(
                            context: context,
                            initialDate: (dueAt ?? DateTime.now()),
                            firstDate: DateTime(2000),
                            lastDate: DateTime(2100),
                          );
                          if (pickedDate == null) return;

                          final pickedTime = await showTimePicker(
                            context: context,
                            initialTime:
                                TimeOfDay.fromDateTime(dueAt ?? DateTime.now()),
                          );
                          if (pickedTime == null) return;

                          final next = DateTime(
                            pickedDate.year,
                            pickedDate.month,
                            pickedDate.day,
                            pickedTime.hour,
                            pickedTime.minute,
                          );

                          setState(() {
                            dueAt = next;
                          });
                        },
                        child: const Text('期限設定'),
                      ),
                      TextButton(
                        onPressed: () {
                          setState(() {
                            dueAt = null;
                          });
                        },
                        child: const Text('クリア'),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  Row(
                    children: [
                      Expanded(
                        child: Text(_formatDateTime(scheduledAt) ?? 'リマインダーなし'),
                      ),
                      TextButton(
                        onPressed: () async {
                          final selection = await showReminderPresetSheet(
                            context,
                            current: scheduledAt,
                          );
                          if (selection == null) return;
                          setState(() {
                            scheduledAt = selection.scheduledAt;
                          });
                        },
                        child: const Text('プリセット'),
                      ),
                    ],
                  ),
                ],
              ),
              actions: [
                TextButton(
                  onPressed: () {
                    Navigator.of(context).pop(null);
                  },
                  child: const Text('キャンセル'),
                ),
                TextButton(
                  onPressed: () {
                    Navigator.of(context).pop(
                      _TaskEditResult(
                        title: controller.text.trim(),
                        dueAt: dueAt,
                        updateDueAt: true,
                        scheduledAt: scheduledAt,
                        updateScheduledAt: true,
                      ),
                    );
                  },
                  child: const Text('保存'),
                ),
              ],
            );
          },
        );
      },
    );

    controller.dispose();
    if (result == null) return;

    if (result.title.isEmpty) {
      showAppErrorSnackBar(context, 'タイトルを入力してください');
      return;
    }

    setState(() {
      _updatingTaskIds.add(task.id);
    });

    try {
      await TasksApi().updateTask(
        id: task.id,
        title: result.title,
        dueAt: result.dueAt,
        updateDueAt: result.updateDueAt,
        scheduledAt: result.scheduledAt,
        updateScheduledAt: result.updateScheduledAt,
      );
      _reload();
    } on ApiException catch (e) {
      if (mounted) {
        showAppErrorSnackBar(context, e.message);
      }
    } catch (_) {
      if (mounted) {
        showAppErrorSnackBar(context, 'タスクの更新に失敗しました');
      }
    } finally {
      if (mounted) {
        setState(() {
          _updatingTaskIds.remove(task.id);
        });
      }
    }
  }

  @override
  void dispose() {
    _titleController.dispose();
    _titleFocusNode.dispose();
    super.dispose();
  }

  void _reload() {
    setState(() {
      _future = TasksApi().fetchTasks();
    });
  }

  Future<void> _createTask() async {
    if (_isCreating) return;

    final title = _titleController.text.trim();
    if (title.isEmpty) {
      showAppErrorSnackBar(context, 'タイトルを入力してください');
      return;
    }

    setState(() {
      _isCreating = true;
    });

    try {
      await TasksApi().createTask(title: title);
      _titleController.clear();
      _titleFocusNode.requestFocus();
      _reload();
    } on ApiException catch (e) {
      if (mounted) {
        showAppErrorSnackBar(context, e.message);
      }
    } catch (_) {
      if (mounted) {
        showAppErrorSnackBar(context, 'タスクの作成に失敗しました');
      }
    } finally {
      if (mounted) {
        setState(() {
          _isCreating = false;
        });
      }
    }
  }

  String _errorText(Object error) {
    if (error is ApiException) {
      return error.message;
    }
    return 'タスクの取得に失敗しました';
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Tasks'),
        actions: [
          IconButton(
            onPressed: _reload,
            icon: const Icon(Icons.refresh),
            tooltip: 'Reload',
          ),
        ],
      ),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _titleController,
                    focusNode: _titleFocusNode,
                    enabled: !_isCreating,
                    textInputAction: TextInputAction.done,
                    onSubmitted: (_) => _createTask(),
                    decoration: const InputDecoration(
                      labelText: 'タスクを追加',
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                ElevatedButton(
                  onPressed: _isCreating ? null : _createTask,
                  child: _isCreating
                      ? const AppLoadingIndicator(size: 18, strokeWidth: 2)
                      : const Text('追加'),
                ),
              ],
            ),
          ),
          const Divider(height: 1),
          Expanded(
            child: FutureBuilder<List<Task>>(
              future: _future,
              builder: (context, snapshot) {
                if (snapshot.connectionState != ConnectionState.done) {
                  return const Center(
                    child: AppLoadingIndicator(size: 24, strokeWidth: 3),
                  );
                }

                if (snapshot.hasError) {
                  return Center(
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Text(_errorText(snapshot.error!)),
                          const SizedBox(height: 12),
                          ElevatedButton(
                            onPressed: _reload,
                            child: const Text('再読み込み'),
                          ),
                        ],
                      ),
                    ),
                  );
                }

                final tasks = snapshot.data ?? const <Task>[];

                if (tasks.isEmpty) {
                  return const Center(child: Text('タスクがありません'));
                }

                return ListView.separated(
                  itemCount: tasks.length,
                  separatorBuilder: (_, __) => const Divider(height: 1),
                  itemBuilder: (context, index) {
                    final task = tasks[index];
                    final isUpdating = task.id.isNotEmpty && _updatingTaskIds.contains(task.id);
                    final dueAtText = _formatDateTime(task.dueAt);
                    final scheduledAtText = _formatDateTime(task.scheduledAt);

                    final subtitleParts = <String>[];
                    if (dueAtText != null) subtitleParts.add('期限: $dueAtText');
                    if (scheduledAtText != null) {
                      subtitleParts.add('リマインダー: $scheduledAtText');
                    }

                    return ListTile(
                      leading: Checkbox(
                        value: task.isDone,
                        onChanged: isUpdating ? null : (_) => _toggleDone(task),
                      ),
                      title: Text(task.title.isEmpty ? '(no title)' : task.title),
                      subtitle: subtitleParts.isEmpty
                          ? null
                          : Text(subtitleParts.join(' / ')),
                      trailing: IconButton(
                        onPressed: isUpdating ? null : () => _editTask(task),
                        icon: const Icon(Icons.edit),
                        tooltip: 'Edit',
                      ),
                    );
                  },
                );
              },
            ),
          ),
        ],
      ),
    );
  }
}

class _TaskEditResult {
  _TaskEditResult({
    required this.title,
    required this.dueAt,
    required this.updateDueAt,
    required this.scheduledAt,
    required this.updateScheduledAt,
  });

  final String title;
  final DateTime? dueAt;
  final bool updateDueAt;
  final DateTime? scheduledAt;
  final bool updateScheduledAt;
}
