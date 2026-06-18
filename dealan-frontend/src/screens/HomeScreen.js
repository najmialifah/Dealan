import React from 'react';
import { View, Text, Button, StyleSheet } from 'react-native';
import AsyncStorage from '@react-native-async-storage/async-storage';

export default function HomeScreen({ navigation }) {
  const handleLogout = async () => {
    await AsyncStorage.removeItem('userToken');
    navigation.replace('Login');
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Selamat Datang di Dealan</Text>
      
      <View style={styles.menuContainer}>
        <Button 
          title="Buat Pesanan (Ride/Car/Send)" 
          onPress={() => navigation.navigate('CreateOrder')} 
        />
      </View>

      <Button title="Logout" color="red" onPress={handleLogout} />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, padding: 20, justifyContent: 'center', alignItems: 'center' },
  title: { fontSize: 24, fontWeight: 'bold', marginBottom: 40 },
  menuContainer: { marginBottom: 30, width: '100%' }
});
