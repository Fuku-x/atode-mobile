import 'package:flutter/material.dart';

import '../../api/api.dart';
import '../../models/task.dart';
import '../../widgets/app_loading_indicator.dart';

import 'tasks_api.dart';

class TasksScreen extends StatefulWidget {
  const TasksScreen({super.key});

  @override
  State<TasksScreen> createState() => _TasksScreenState();
}

class _TasksScreenState extends State<TasksScreen> {
  late Future<List<Task>> _future;

  @override
  void initState() {
    super.initState();
    _future = TasksApi().fetchTasks();
  }

  void _reload() {
    setState(() {
      _future = TasksApi().fetchTasks();
    });
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
      body: FutureBuilder<List<Task>>(
        future: _future,
        builder: (context, snapshot) {
          if (snapshot.connectionState != ConnectionState.done) {
            return const Center(child: AppLoadingIndicator(size: 24, strokeWidth: 3));
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
                leading: Icon(task.isDone ? Icons.check_circle : Icons.radio_button_unchecked),
                title: Text(task.title.isEmpty ? '(no title)' : task.title),
                subtitle: task.id.isEmpty ? null : Text(task.id),
              );
            },
          );
        },
      ),
    );
  }
}
