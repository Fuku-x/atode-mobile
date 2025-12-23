class Task {
  Task({
    required this.id,
    required this.title,
    required this.isDone,
    this.dueAt,
  });

  final String id;
  final String title;
  final bool isDone;
  final DateTime? dueAt;

  factory Task.fromJson(Map<String, dynamic> json) {
    final title = (json['title'] ?? json['name'] ?? '').toString();
    final id = (json['id'] ?? json['taskId'] ?? '').toString();
    final doneRaw = json['isDone'] ?? json['done'] ?? json['completed'] ?? false;

    final dueRaw = json['dueAt'] ?? json['due_at'] ?? json['due_at_ms'];
    DateTime? dueAt;
    if (dueRaw is String && dueRaw.isNotEmpty) {
      dueAt = DateTime.tryParse(dueRaw);
    } else if (dueRaw is int) {
      dueAt = DateTime.fromMillisecondsSinceEpoch(dueRaw);
    }

    return Task(
      id: id,
      title: title,
      isDone: doneRaw == true,
      dueAt: dueAt,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'isDone': isDone,
      if (dueAt != null) 'dueAt': dueAt!.toIso8601String(),
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
