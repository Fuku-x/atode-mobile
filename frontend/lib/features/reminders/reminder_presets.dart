enum ReminderPreset {
  in1Hour,
  tonight,
  tomorrow,
}

DateTime computeReminderDateTime(ReminderPreset preset, DateTime now) {
  switch (preset) {
    case ReminderPreset.in1Hour:
      return now.add(const Duration(hours: 1));
    case ReminderPreset.tonight:
      final base = DateTime(now.year, now.month, now.day, 21);
      if (base.isAfter(now)) return base;
      return base.add(const Duration(days: 1));
    case ReminderPreset.tomorrow:
      return DateTime(now.year, now.month, now.day + 1, 9);
  }
}

String reminderPresetLabel(ReminderPreset preset) {
  switch (preset) {
    case ReminderPreset.in1Hour:
      return '1時間後';
    case ReminderPreset.tonight:
      return '今夜';
    case ReminderPreset.tomorrow:
      return '明日';
  }
}
