import '../auth/firebase_id_token_provider.dart';

import 'api_client.dart';
import 'api_config.dart';

class FirebaseApiClient {
  static ApiClient create({
    ApiConfig? config,
    bool forceRefreshToken = false,
  }) {
    final configInstance = config ?? ApiConfig.fromEnvironment();

    return ApiClient(
      config: configInstance,
      tokenProvider: () => getFirebaseIdToken(forceRefresh: forceRefreshToken),
    );
  }
}
