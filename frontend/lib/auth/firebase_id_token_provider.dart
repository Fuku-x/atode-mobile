import 'package:firebase_auth/firebase_auth.dart';

Future<String?> getFirebaseIdToken({
  FirebaseAuth? auth,
  bool forceRefresh = false,
}) async {
  final a = auth ?? FirebaseAuth.instance;
  final user = a.currentUser;
  if (user == null) return null;
  return user.getIdToken(forceRefresh);
}

Future<String?> firebaseIdTokenProvider() {
  return getFirebaseIdToken();
}
