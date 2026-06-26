import 'dart:io';

import 'package:drift/drift.dart';
import 'package:drift/native.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:path/path.dart' as p;
import 'package:path_provider/path_provider.dart';

part 'schema.g.dart';

class SyncQueue extends Table {
  TextColumn get id => text()();
  TextColumn get eventType => text().named('event_type')();
  TextColumn get payload => text()();
  TextColumn get clientTimestamp => text().named('client_timestamp')();
  TextColumn get status => text().withDefault(const Constant('pending'))();
  IntColumn get retryCount => integer().named('retry_count').withDefault(const Constant(0))();
  TextColumn get createdAt => text().named('created_at')();

  @override
  Set<Column> get primaryKey => {id};
}

class Operations extends Table {
  TextColumn get id => text()();
  TextColumn get plotId => text().named('plot_id')();
  TextColumn get type => text()();
  TextColumn get date => text()();
  TextColumn get responsible => text().withDefault(const Constant(''))();
  TextColumn get productUsed => text().named('product_used').withDefault(const Constant(''))();
  RealColumn get quantity => real().withDefault(const Constant(0))();
  RealColumn get cost => real().withDefault(const Constant(0))();
  TextColumn get notes => text().withDefault(const Constant(''))();
  IntColumn get synced => integer().withDefault(const Constant(0))();
  TextColumn get createdAt => text().named('created_at')();

  @override
  Set<Column> get primaryKey => {id};
}

class Plots extends Table {
  TextColumn get id => text()();
  TextColumn get farmId => text().named('farm_id')();
  TextColumn get name => text()();
  RealColumn get areaHa => real().named('area_ha').withDefault(const Constant(0))();
  TextColumn get cultivar => text().withDefault(const Constant(''))();
  IntColumn get synced => integer().withDefault(const Constant(0))();

  @override
  Set<Column> get primaryKey => {id};
}

class Farms extends Table {
  TextColumn get id => text()();
  TextColumn get name => text()();
  TextColumn get owner => text().withDefault(const Constant(''))();
  TextColumn get location => text().withDefault(const Constant(''))();
  IntColumn get synced => integer().withDefault(const Constant(0))();

  @override
  Set<Column> get primaryKey => {id};
}

@DriftDatabase(tables: [SyncQueue, Operations, Plots, Farms])
class AppDatabase extends _$AppDatabase {
  AppDatabase() : super(_openConnection());

  @override
  int get schemaVersion => 1;

  static QueryExecutor _openConnection() {
    return LazyDatabase(() async {
      final dir = await getApplicationDocumentsDirectory();
      final file = File(p.join(dir.path, 'cafeos_offline.db'));
      return NativeDatabase(file);
    });
  }
}

final databaseProvider = Provider<AppDatabase>((ref) => AppDatabase());
