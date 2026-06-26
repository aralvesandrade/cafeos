import 'dart:convert';
import 'package:drift/drift.dart';
import 'package:uuid/uuid.dart';
import '../db/schema.dart' as db;

class SyncQueueItem {
  final String id;
  final String eventType;
  final String payload;
  final String clientTimestamp;
  final String status;
  final int retryCount;
  final String createdAt;

  SyncQueueItem({
    required this.id,
    required this.eventType,
    required this.payload,
    required this.clientTimestamp,
    required this.status,
    this.retryCount = 0,
    this.createdAt = '',
  });
}

class SyncQueueRepo {
  final db.AppDatabase _db;

  SyncQueueRepo(this._db);

  Future<List<SyncQueueItem>> getPendingItems() async {
    final rows = await (_db.select(_db.syncQueue)
          ..where((t) => t.status.equals('pending'))
          ..orderBy([(t) => OrderingTerm(expression: t.createdAt)])
          ..limit(50))
        .get();
    return rows.map(_fromRow).toList();
  }

  Future<List<SyncQueueItem>> getAll() async {
    final rows = await (_db.select(_db.syncQueue)
          ..orderBy([(t) => OrderingTerm(expression: t.createdAt, mode: OrderingMode.desc)]))
        .get();
    return rows.map(_fromRow).toList();
  }

  Future<void> enqueue(String eventType, Map<String, dynamic> payload) async {
    final id = const Uuid().v4();
    final now = DateTime.now().toUtc().toIso8601String();
    await _db.into(_db.syncQueue).insert(db.SyncQueueCompanion.insert(
          id: id,
          eventType: eventType,
          payload: jsonEncode(payload),
          clientTimestamp: now,
          createdAt: now,
        ));
  }

  Future<void> markSynced(String id) async {
    await (_db.update(_db.syncQueue)
          ..where((t) => t.id.equals(id)))
        .write(const db.SyncQueueCompanion(status: Value('synced')));
  }

  Future<void> markFailed(String id) async {
    await (_db.update(_db.syncQueue)
          ..where((t) => t.id.equals(id)))
        .write(const db.SyncQueueCompanion(status: Value('failed')));
  }

  Future<void> incrementRetry(String id) async {
    final row = await (_db.select(_db.syncQueue)..where((t) => t.id.equals(id))).getSingle();
    final newCount = row.retryCount + 1;
    if (newCount >= 3) {
      await markFailed(id);
    } else {
      await (_db.update(_db.syncQueue)
            ..where((t) => t.id.equals(id)))
          .write(db.SyncQueueCompanion(
            retryCount: Value(newCount),
            status: const Value('pending'),
          ));
    }
  }

  SyncQueueItem _fromRow(db.SyncQueueData row) => SyncQueueItem(
        id: row.id,
        eventType: row.eventType,
        payload: row.payload,
        clientTimestamp: row.clientTimestamp,
        status: row.status,
        retryCount: row.retryCount,
        createdAt: row.createdAt,
      );
}
