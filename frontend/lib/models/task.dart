class Task {
  Task({
    required this.id,
    required this.title,
    required this.isDone,
    this.dueAt,
    this.scheduledAt,
  });

  final String id;
  final String title;
  final bool isDone;
  final DateTime? dueAt;
  final DateTime? scheduledAt;

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

    final scheduledRaw =
        json['scheduledAt'] ?? json['scheduled_at'] ?? json['scheduled_at_ms'];
    DateTime? scheduledAt;
    if (scheduledRaw is String && scheduledRaw.isNotEmpty) {
      scheduledAt = DateTime.tryParse(scheduledRaw);
    } else if (scheduledRaw is int) {
      scheduledAt = DateTime.fromMillisecondsSinceEpoch(scheduledRaw);
    }

    return Task(
      id: id,
      title: title,
      isDone: doneRaw == true,
      dueAt: dueAt,
      scheduledAt: scheduledAt,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'isDone': isDone,
      if (dueAt != null) 'dueAt': dueAt!.toIso8601String(),
      if (scheduledAt != null) 'scheduledAt': scheduledAt!.toIso8601String(),
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
