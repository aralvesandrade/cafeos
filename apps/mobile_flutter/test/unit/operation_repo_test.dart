import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_flutter/models/models.dart';

void main() {
  group('Operation', () {
    test('fromJson parses correctly', () {
      final json = {
        'id': '123',
        'plot_id': 'plot-1',
        'type': 'adubacao',
        'date': '2026-06-25',
        'responsible': 'João',
        'product_used': 'NPK',
        'quantity': 50.0,
        'cost': 2500.0,
        'notes': 'Manual',
        'synced': 0,
      };

      final op = Operation.fromJson(json);

      expect(op.id, '123');
      expect(op.plotId, 'plot-1');
      expect(op.type, 'adubacao');
      expect(op.quantity, 50.0);
      expect(op.cost, 2500.0);
    });

    test('toJson produces correct output', () {
      final op = Operation(
        id: '123',
        plotId: 'plot-1',
        type: 'colheita',
        date: '2026-06-25',
        cost: 1000.0,
      );

      final json = op.toJson();

      expect(json['id'], '123');
      expect(json['type'], 'colheita');
      expect(json['cost'], 1000.0);
    });

    test('copyWith preserves unchanged fields', () {
      final op = Operation(id: '1', plotId: 'p1', type: 'poda', date: '2026-01-01');
      final copy = op.copyWith(type: 'irrigacao');

      expect(copy.id, '1');
      expect(copy.plotId, 'p1');
      expect(copy.type, 'irrigacao');
      expect(copy.date, '2026-01-01');
    });
  });

  group('Plot', () {
    test('fromJson parses correctly', () {
      final json = {
        'id': 'plot-1',
        'farm_id': 'farm-1',
        'name': 'Talhão A',
        'area_ha': 5.5,
        'cultivar': 'Catuaí',
      };

      final plot = Plot.fromJson(json);

      expect(plot.id, 'plot-1');
      expect(plot.farmId, 'farm-1');
      expect(plot.name, 'Talhão A');
      expect(plot.areaHa, 5.5);
    });
  });
}
