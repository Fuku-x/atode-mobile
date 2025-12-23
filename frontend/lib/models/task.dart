class Task {
  Task({
    required this.id,
    required this.title,
    required this.isDone,
  });

  final String id;
  final String title;
  final bool isDone;

  factory Task.fromJson(Map<String, dynamic> json) {
    final title = (json['title'] ?? json['name'] ?? '').toString();
    final id = (json['id'] ?? json['taskId'] ?? '').toString();
    final doneRaw = json['isDone'] ?? json['done'] ?? json['completed'] ?? false;

    return Task(
      id: id,
      title: title,
      isDone: doneRaw == true,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'isDone': isDone,
    };
  }

  static List<Task> listFromJson(dynamic json) {
    if (json is List) {
      return json
          .whereType<Map>()
          .map((e) => Task.fromJson(Map<String, dynamic>.from(e)))
          .toList();
    }

    if (json is Map) {
      final tasks = json['tasks'];
      if (tasks is List) {
        return tasks
            .whereType<Map>()
            .map((e) => Task.fromJson(Map<String, dynamic>.from(e)))
            .toList();
      }
    }

    return const [];
  }
}
