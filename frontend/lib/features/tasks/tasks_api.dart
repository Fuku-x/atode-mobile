import '../../api/api.dart';
import '../../models/task.dart';

class TasksApi {
  TasksApi({ApiClient? client}) : _client = client ?? FirebaseApiClient.create();

  final ApiClient _client;

  Future<List<Task>> fetchTasks() {
    return _client.get<List<Task>>(
      '/tasks',
      decoder: (json) => Task.listFromJson(json),
    );
  }

  Future<void> createTask({required String title}) async {
    await _client.post<dynamic>(
      '/tasks',
      body: {
        'title': title,
      },
    );
  }

  Future<void> updateTask({
    required String id,
    String? title,
    bool? isDone,
    DateTime? dueAt,
    bool updateDueAt = false,
    DateTime? scheduledAt,
    bool updateScheduledAt = false,
  }) async {
    final body = <String, dynamic>{};
    if (title != null) body['title'] = title;
    if (isDone != null) body['isDone'] = isDone;
    if (updateDueAt) {
      body['dueAt'] = dueAt?.toIso8601String();
    }
    if (updateScheduledAt) {
      body['scheduledAt'] = scheduledAt?.toIso8601String();
    }

    if (body.isEmpty) return;

    await _client.put<dynamic>(
      '/tasks/$id',
      body: body,
    );
  }
}
