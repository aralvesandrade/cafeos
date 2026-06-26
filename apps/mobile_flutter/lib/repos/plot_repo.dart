import '../db/schema.dart';
import '../models/models.dart';

class PlotRepo {
  final AppDatabase _db;

  PlotRepo(this._db);

  Future<List<Plot>> getAll() async {
    final rows = await (_db.select(_db.plots)
          ..orderBy([(t) => OrderingTerm(expression: t.name)]))
        .get();
    return rows.map(_fromRow).toList();
  }

  Future<void> upsertAll(List<Plot> plots) async {
    for (final p in plots) {
      await _db.into(_db.plots).insertOnConflictUpdate(PlotsCompanion.insert(
            id: Value(p.id),
            farmId: Value(p.farmId),
            name: Value(p.name),
            areaHa: Value(p.areaHa),
            cultivar: Value(p.cultivar),
          ));
    }
  }

  Plot _fromRow(PlotsData row) => Plot(
        id: row.id,
        farmId: row.farmId,
        name: row.name,
        areaHa: row.areaHa,
        cultivar: row.cultivar,
        synced: row.synced,
      );
}
