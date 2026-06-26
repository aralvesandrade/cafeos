import 'package:flutter/material.dart';
import '../shared/theme/app_theme.dart';

class StatusBadge extends StatelessWidget {
  final String status;
  final double fontSize;

  const StatusBadge({
    super.key,
    required this.status,
    this.fontSize = 11,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
      decoration: BoxDecoration(
        color: backgroundColor.withValues(alpha: 0.15),
        borderRadius: BorderRadius.circular(4),
      ),
      child: Text(
        label,
        style: TextStyle(fontSize: fontSize, color: backgroundColor, fontWeight: FontWeight.w600),
      ),
    );
  }

  Color get backgroundColor {
    switch (status) {
      case 'pending':
        return Colors.amber;
      case 'synced':
        return Colors.green;
      case 'failed':
        return Colors.red;
      default:
        return AppTheme.textSecondary;
    }
  }

  String get label {
    switch (status) {
      case 'pending':
        return 'Pendente';
      case 'synced':
        return 'Sincronizado';
      case 'failed':
        return 'Falhou';
      default:
        return status;
    }
  }
}

class SyncButton extends StatelessWidget {
  final bool syncing;
  final VoidCallback? onPressed;

  const SyncButton({super.key, this.syncing = false, this.onPressed});

  @override
  Widget build(BuildContext context) {
    return IconButton(
      icon: syncing
          ? const SizedBox(
              width: 20,
              height: 20,
              child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
            )
          : const Icon(Icons.sync),
      onPressed: syncing ? null : onPressed,
    );
  }
}
