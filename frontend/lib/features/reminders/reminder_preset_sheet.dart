import 'package:flutter/material.dart';

import 'reminder_presets.dart';

class ReminderSelection {
  ReminderSelection({required this.scheduledAt, required this.updateScheduledAt});

  final DateTime? scheduledAt;
  final bool updateScheduledAt;
}

Future<ReminderSelection?> showReminderPresetSheet(
  BuildContext context, {
  DateTime? current,
}) {
  return showModalBottomSheet<ReminderSelection>(
    context: context,
    showDragHandle: true,
    builder: (context) {
      return SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              title: const Text('リマインダー'),
              subtitle: Text(current == null ? '未設定' : current.toLocal().toString()),
            ),
            const Divider(height: 1),
            ...ReminderPreset.values.map(
              (preset) => ListTile(
                title: Text(reminderPresetLabel(preset)),
                onTap: () {
                  final dt = computeReminderDateTime(preset, DateTime.now());
                  Navigator.of(context).pop(
                    ReminderSelection(scheduledAt: dt, updateScheduledAt: true),
                  );
                },
              ),
            ),
            ListTile(
              title: const Text('クリア'),
              onTap: () {
                Navigator.of(context).pop(
                  ReminderSelection(scheduledAt: null, updateScheduledAt: true),
                );
              },
            ),
            const SizedBox(height: 12),
          ],
        ),
      );
    },
  );
}
