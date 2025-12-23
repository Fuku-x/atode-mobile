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
}
