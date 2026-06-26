import 'package:flutter/material.dart';
import 'router/app_router.dart';
import 'shared/theme/app_theme.dart';

class CafeOSApp extends StatelessWidget {
  const CafeOSApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      title: 'CafeOS',
      theme: AppTheme.light,
      routerConfig: appRouter,
      debugShowCheckedModeBanner: false,
    );
  }
}
