import 'package:flutter/material.dart';

import 'app_route_names.dart';
import '../screens/home_screen.dart';
import '../screens/task_screen.dart';

class AppRouter {
  static Route<dynamic> onGenerateRoute(RouteSettings settings) {
    switch (settings.name) {
      case AppRouteNames.home:
        return MaterialPageRoute<void>(
          settings: settings,
          builder: (_) => const HomeScreen(),
        );
      case AppRouteNames.tasks:
        return MaterialPageRoute<void>(
          settings: settings,
          builder: (_) => const TaskScreen(),
        );
      default:
        return MaterialPageRoute<void>(
          settings: settings,
          builder: (_) => Scaffold(
            appBar: AppBar(title: const Text('Not Found')),
            body: Center(child: Text('Route not found: ${settings.name}')),
          ),
        );
    }
  }
}
