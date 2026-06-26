import 'dart:convert';
import '../api/client.dart';
import '../api/storage.dart';
import '../repos/sync_queue_repo.dart';

class SyncService {
  final SyncQueueRepo _repo;

  SyncService(this._repo);

  Future<SyncResult> syncAll() async {
    final token = await Storage.getToken();
    if (token == null) return SyncResult(synced: 0, failed: 0);

    final items = await _repo.getPendingItems();
    if (items.isEmpty) return SyncResult(synced: 0, failed: 0);

    int synced = 0;
    int failed = 0;

    final batch = items.map((item) => {
          'client_id': item.id,
          'event_type': item.eventType,
          'payload': _parsePayload(item.payload),
          'client_timestamp': item.clientTimestamp,
        }).toList();

    try {
      final response = await ApiClient.request(
        '/sync',
        method: 'POST',
        body: {'batch': batch},
      );

      if (response.statusCode == 202) {
        final data = response.data as Map<String, dynamic>;
        final errors = data['errors'] as List<dynamic>? ?? [];
        final errorIds = errors.map((e) => (e as Map)['client_id'] as String).toSet();

        for (final item in items) {
          if (errorIds.contains(item.id)) {
            await _repo.incrementRetry(item.id);
            failed++;
          } else {
            await _repo.markSynced(item.id);
            synced++;
          }
        }
      }
    } catch (e) {
      for (final item in items) {
        await _repo.incrementRetry(item.id);
      }
      failed = items.length;
    }

    return SyncResult(synced: synced, failed: failed);
  }

  Map<String, dynamic> _parsePayload(String payload) {
    try {
      return jsonDecode(payload) as Map<String, dynamic>;
    } catch (_) {
      return {};
    }
  }
}

class SyncResult {
  final int synced;
  final int failed;
  SyncResult({required this.synced, required this.failed});
}
