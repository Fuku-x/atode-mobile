import 'package:flutter/material.dart';

import 'app_loading_indicator.dart';

class LoadingOutlinedButton extends StatelessWidget {
  const LoadingOutlinedButton({
    super.key,
    required this.isLoading,
    required this.onPressed,
    required this.child,
  });

  final bool isLoading;
  final VoidCallback? onPressed;
  final Widget child;

  @override
  Widget build(BuildContext context) {
    return OutlinedButton(
      onPressed: isLoading ? null : onPressed,
      child: isLoading ? const AppLoadingIndicator() : child,
    );
  }
}
