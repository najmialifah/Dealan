import React, { useEffect, useState } from 'react';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { ActivityIndicator, View } from 'react-native';

// Screens
import LoginScreen from '../screens/LoginScreen';
import RegisterScreen from '../screens/RegisterScreen';
import HomeScreen from '../screens/HomeScreen';
import CreateOrderScreen from '../screens/CreateOrderScreen';
import NegotiationScreen from '../screens/NegotiationScreen';
import MatchingScreen from '../screens/MatchingScreen';
import PaymentScreen from '../screens/PaymentScreen';
import RatingScreen from '../screens/RatingScreen';

const Stack = createNativeStackNavigator();

export default function AppNavigator() {
  const [isLoading, setIsLoading] = useState(true);
  const [userToken, setUserToken] = useState(null);

  useEffect(() => {
    const checkToken = async () => {
      try {
        const token = await AsyncStorage.getItem('userToken');
        setUserToken(token);
      } catch (e) {
        console.error(e);
      } finally {
        setIsLoading(false);
      }
    };
    checkToken();
  }, []);

  if (isLoading) {
    return (
      <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center' }}>
        <ActivityIndicator size="large" color="#0000ff" />
      </View>
    );
  }

  return (
    <Stack.Navigator>
      {userToken == null ? (
        // Auth Stack
        <>
          <Stack.Screen name="Login" component={LoginScreen} />
          <Stack.Screen name="Register" component={RegisterScreen} />
        </>
      ) : (
        // Main Stack
        <>
          <Stack.Screen name="Home" component={HomeScreen} />
          <Stack.Screen name="CreateOrder" component={CreateOrderScreen} />
          <Stack.Screen name="Negotiation" component={NegotiationScreen} />
          <Stack.Screen name="Matching" component={MatchingScreen} />
          <Stack.Screen name="Payment" component={PaymentScreen} />
          <Stack.Screen name="Rating" component={RatingScreen} />
        </>
      )}
    </Stack.Navigator>
  );
}
