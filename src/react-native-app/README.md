# Northstar Commerce mobile app

This Expo app gives operators and engineers a mobile client for the Northstar Commerce storefront. It shares the same backend APIs as the web storefront while exercising mobile networking, tracing, cart, and checkout flows.

## Get started

Start the Northstar Commerce stack from the repository root:

```bash
make start # or start-minimal
```

Install the mobile dependencies:

```bash
cd src/react-native-app
npm install
```

## Run on Android

The Android target builds the app, deploys it to a running emulator or connected device, and starts the JavaScript bundle server:

```bash
npm run android
```

## Run on iOS

Install CocoaPods dependencies the first time you build the iOS app:

```bash
cd ios && pod install && cd ..
```

Start the bundle server:

```bash
npm run start
```

Then open `src/react-native-app/ios/react-native-app.xcworkspace` in Xcode and run the app. You can also build from the command line:

```bash
npm run ios
```

## Build Android in a container

If the Android toolchain is not installed on your host, build the APK with Docker from the repository root:

```bash
make build-react-native-android
```

Or build directly from this folder:

```bash
docker build -f android.Dockerfile --platform=linux/amd64 --output=. .
```

The build writes `react-native-app.apk` to the current directory.

## Pointing to another environment

By default, the app connects to `EXPO_PUBLIC_FRONTEND_PROXY_PORT` on localhost. Use the Settings tab to point the app at another Northstar Commerce environment.

## Troubleshooting

For JavaScript bundle issues, remove `src/react-native-app/node_modules/` and run `npm install` again.

For Android build issues, stop Gradle daemons and clear local caches:

```bash
cd src/react-native-app/android
./gradlew --stop
rm -rf ~/.gradle/caches/
```

For iOS dependency drift, install the pinned Ruby dependencies before running CocoaPods:

```bash
rbenv install 2.7.6
bundle install
cd ios
bundle exec pod install
```

If Xcode continues to fail after pod installation, clear derived data and retry the build:

```bash
rm -rf ~/Library/Developer/Xcode/DerivedData
```
