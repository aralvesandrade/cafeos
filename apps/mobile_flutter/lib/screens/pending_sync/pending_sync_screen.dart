import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../services/services.dart';
import '../../services/sync_service.dart';
import '../../shared/theme/app_theme.dart';

final pendingSyncControllerProvider =
    StateNotifierProvider<PendingSyncController, PendingSyncState>((ref) {
  return PendingSyncController(
    ref.read(syncServiceProvider),
    ref.read(offlineServiceProvider),
  );
});

class PendingSyncState {
  final List<SyncItemView> items;
  final bool loading;
  final bool syncing;

  PendingSyncState({
    this.items = const [],
    this.loading = true,
    this.syncing = false,
  });

  PendingSyncState copyWith({
    List<SyncItemView>? items,
    bool? loading,
    bool? syncing,
  }) =>
      PendingSyncState(
        items: items ?? this.items,
        loading: loading ?? this.loading,
        syncing: syncing ?? this.syncing,
      );
}

class SyncItemView {
  final String id;
  final String eventType;
  final String status;
  final int retryCount;
  final String createdAt;

  SyncItemView({
    required this.id,
    required this.eventType,
    required this.status,
    this.retryCount = 0,
    this.createdAt = '',
  });
}

class PendingSyncController extends StateNotifier<PendingSyncState> {
  final SyncService _syncService;
  final OfflineService _offlineService;

  PendingSyncController(this._syncService, this._offlineService)
      : super(PendingSyncState()) {
    load();
  }

  Future<void> load() async {
    state = state.copyWith(loading: true);
    try {
      final items = await _offlineService.syncQueueRepo.getAll();
      state = state.copyWith(
        items: items
            .map((i) => SyncItemView(
                  id: i.id,
                  eventType: i.eventType,
                  status: i.status,
                  retryCount: i.retryCount,
                  createdAt: i.createdAt,
                ))
            .toList(),
        loading: false,
      );
    } catch (_) {
      state = state.copyWith(loading: false);
    }
  }

  Future<void> syncAll() async {
    state = state.copyWith(syncing: true);
    await _syncService.syncAll();
    state = state.copyWith(syncing: false);
    await load();
  }
}

class PendingSyncScreen extends ConsumerWidget {
  const PendingSyncScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(pendingSyncControllerProvider);
    final controller = ref.read(pendingSyncControllerProvider.notifier);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Pendências'),
        actions: [
          IconButton(
            icon: state.syncing
                ? const SizedBox(
                    width: 20,
                    height: 20,
                    child: CircularProgressIndicator(
                      strokeWidth: 2,
                      color: Colors.white,
                    ),
                  )
                : const Icon(Icons.sync),
            onPressed: state.syncing ? null : controller.syncAll,
          ),
        ],
      ),
      body: _buildBody(state),
    );
  }

  Widget _buildBody(PendingSyncState state) {
    if (state.loading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (state.items.isEmpty) {
      return const Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.check_circle_outline, size: 64, color: AppTheme.green),
            SizedBox(height: 16),
            Text(
              'Nenhum item pendente',
              style: TextStyle(color: AppTheme.textSecondary, fontSize: 16),
            ),
          ],
        ),
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: state.items.length,
      itemBuilder: (context, index) {
        final item = state.items[index];
        return Card(
          margin: const EdgeInsets.only(bottom: 8),
          child: ListTile(
            leading: _statusIcon(item.status),
            title: Text(item.eventType),
            subtitle: Text('${item.createdAt}\nTentativas: ${item.retryCount}'),
            trailing: Text(
              _statusLabel(item.status),
              style: TextStyle(
                color: _statusColor(item.status),
                fontWeight: FontWeight.w600,
                fontSize: 12,
              ),
            ),
          ),
        );
      },
    );
  }

  Widget _statusIcon(String status) {
    switch (status) {
      case 'pending':
        return const Icon(Icons.hourglass_empty, color: Colors.amber);
      case 'synced':
        return const Icon(Icons.check_circle, color: Colors.green);
      case 'failed':
        return const Icon(Icons.error, color: Colors.red);
      default:
        return const Icon(Icons.help_outline);
    }
  }

  Color _statusColor(String status) {
    switch (status) {
      case 'pending':
        return Colors.amber.shade700;
      case 'synced':
        return Colors.green;
      case 'failed':
        return Colors.red;
      default:
        return Colors.grey;
    }
  }

  String _statusLabel(String status) {
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
