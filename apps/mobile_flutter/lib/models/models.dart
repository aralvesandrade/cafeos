class Operation {
  final String id;
  final String plotId;
  final String type;
  final String date;
  final String responsible;
  final String productUsed;
  final double quantity;
  final double cost;
  final String notes;
  final int synced;
  final String createdAt;

  Operation({
    required this.id,
    required this.plotId,
    required this.type,
    required this.date,
    this.responsible = '',
    this.productUsed = '',
    this.quantity = 0,
    this.cost = 0,
    this.notes = '',
    this.synced = 0,
    this.createdAt = '',
  });

  Map<String, dynamic> toJson() => {
        'id': id,
        'plot_id': plotId,
        'type': type,
        'date': date,
        'responsible': responsible,
        'product_used': productUsed,
        'quantity': quantity,
        'cost': cost,
        'notes': notes,
        'synced': synced,
        'created_at': createdAt,
      };

  factory Operation.fromJson(Map<String, dynamic> json) => Operation(
        id: json['id'] as String,
        plotId: json['plot_id'] as String,
        type: json['type'] as String,
        date: json['date'] as String,
        responsible: json['responsible'] as String? ?? '',
        productUsed: json['product_used'] as String? ?? '',
        quantity: (json['quantity'] as num?)?.toDouble() ?? 0,
        cost: (json['cost'] as num?)?.toDouble() ?? 0,
        notes: json['notes'] as String? ?? '',
        synced: json['synced'] as int? ?? 0,
        createdAt: json['created_at'] as String? ?? '',
      );

  Operation copyWith({
    String? id,
    String? plotId,
    String? type,
    String? date,
    String? responsible,
    String? productUsed,
    double? quantity,
    double? cost,
    String? notes,
    int? synced,
    String? createdAt,
  }) =>
      Operation(
        id: id ?? this.id,
        plotId: plotId ?? this.plotId,
        type: type ?? this.type,
        date: date ?? this.date,
        responsible: responsible ?? this.responsible,
        productUsed: productUsed ?? this.productUsed,
        quantity: quantity ?? this.quantity,
        cost: cost ?? this.cost,
        notes: notes ?? this.notes,
        synced: synced ?? this.synced,
        createdAt: createdAt ?? this.createdAt,
      );
}

class Plot {
  final String id;
  final String farmId;
  final String name;
  final double areaHa;
  final String cultivar;
  final int synced;

  Plot({
    required this.id,
    required this.farmId,
    required this.name,
    this.areaHa = 0,
    this.cultivar = '',
    this.synced = 0,
  });

  factory Plot.fromJson(Map<String, dynamic> json) => Plot(
        id: json['id'] as String,
        farmId: json['farm_id'] as String,
        name: json['name'] as String,
        areaHa: (json['area_ha'] as num?)?.toDouble() ?? 0,
        cultivar: json['cultivar'] as String? ?? '',
        synced: json['synced'] as int? ?? 0,
      );
}

class Farm {
  final String id;
  final String name;
  final String owner;
  final String location;
  final int synced;

  Farm({
    required this.id,
    required this.name,
    this.owner = '',
    this.location = '',
    this.synced = 0,
  });

  factory Farm.fromJson(Map<String, dynamic> json) => Farm(
        id: json['id'] as String,
        name: json['name'] as String,
        owner: json['owner'] as String? ?? '',
        location: json['location'] as String? ?? '',
        synced: json['synced'] as int? ?? 0,
      );
}
