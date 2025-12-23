import 'dart:developer';

import 'package:firebase_auth/firebase_auth.dart';

Future<String?> getFirebaseIdToken({
  FirebaseAuth? auth,
  bool forceRefresh = false,
}) async {
  final authInstance = auth ?? FirebaseAuth.instance;
  final user = authInstance.currentUser;
  if (user == null) return null;

  try {
    return await user.getIdToken(forceRefresh);
  } catch (e, st) {
    log(
      'Failed to get Firebase ID token',
      error: e,
      stackTrace: st,
      name: 'firebase_id_token_provider',
    );
    return null;
  }
}

Future<String?> firebaseIdTokenProvider() {
  return getFirebaseIdToken();
}
