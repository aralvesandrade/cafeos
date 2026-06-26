import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_flutter/app.dart';

void main() {
  testWidgets('App renders login screen', (tester) async {
    await tester.pumpWidget(const ProviderScope(child: CafeOSApp()));
    await tester.pumpAndSettle();

    expect(find.text('CafeOS'), findsOneWidget);
    expect(find.text('Entrar'), findsOneWidget);
  });
}
