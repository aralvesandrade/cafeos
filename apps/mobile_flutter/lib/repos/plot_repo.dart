import 'package:drift/drift.dart';
import '../db/schema.dart' as db;
import '../models/models.dart' as models;

class PlotRepo {
  final db.AppDatabase _db;

  PlotRepo(this._db);

  Future<List<models.Plot>> getAll() async {
    final rows = await (_db.select(_db.plots)
          ..orderBy([(t) => OrderingTerm(expression: t.name)]))
        .get();
    return rows.map(_fromRow).toList();
  }

  Future<void> upsertAll(List<models.Plot> plots) async {
    for (final p in plots) {
      await _db.into(_db.plots).insertOnConflictUpdate(db.PlotsCompanion.insert(
            id: p.id,
            farmId: p.farmId,
            name: p.name,
            areaHa: Value(p.areaHa),
            cultivar: Value(p.cultivar),
          ));
    }
  }

  models.Plot _fromRow(db.Plot row) => models.Plot(
        id: row.id,
        farmId: row.farmId,
        name: row.name,
        areaHa: row.areaHa,
        cultivar: row.cultivar,
        synced: row.synced,
      );
}
