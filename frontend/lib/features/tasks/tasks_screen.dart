import 'package:flutter/material.dart';

import '../../api/api.dart';
import '../../models/task.dart';
import '../../widgets/app_loading_indicator.dart';
import '../../widgets/app_snack_bar.dart';

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

  @override
  void initState() {
    super.initState();
    _future = TasksApi().fetchTasks();
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
                    return ListTile(
                      leading: Icon(
                        task.isDone
                            ? Icons.check_circle
                            : Icons.radio_button_unchecked,
                      ),
                      title: Text(task.title.isEmpty ? '(no title)' : task.title),
                      subtitle: task.id.isEmpty ? null : Text(task.id),
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
