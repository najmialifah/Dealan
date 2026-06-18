import React, { useState } from 'react';
import { View, Text, TextInput, Button, StyleSheet, ActivityIndicator, Alert } from 'react-native';
import { submitRating } from '../services/ratingApi';

export default function RatingScreen({ route, navigation }) {
  const { order_id, driver_id } = route.params || {};
  const [score, setScore] = useState(5);
  const [comment, setComment] = useState('');
  const [loading, setLoading] = useState(false);

  const handleRating = async () => {
    try {
      setLoading(true);
      const payload = {
        order_id: String(order_id),
        reviewer_id: 'usr-1', // Mock
        reviewer_role: 'user',
        target_id: String(driver_id),
        target_role: 'driver',
        rating_score: score,
        comment
      };

      await submitRating(payload);
      Alert.alert('Terima Kasih', 'Rating telah dikirimkan', [
        { text: 'Kembali ke Home', onPress: () => navigation.popToTop() }
      ]);
    } catch (err) {
      // Handled globally
    } finally {
      setLoading(false);
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Beri Penilaian Driver</Text>
      
      <Text style={styles.label}>Skor (1-5):</Text>
      <View style={styles.scoreContainer}>
        {[1, 2, 3, 4, 5].map((val) => (
          <Button 
            key={val} 
            title={String(val)} 
            color={score === val ? 'orange' : 'gray'} 
            onPress={() => setScore(val)} 
          />
        ))}
      </View>

      <Text style={styles.label}>Komentar:</Text>
      <TextInput
        style={styles.input}
        multiline
        numberOfLines={4}
        placeholder="Bagaimana perjalanan Anda?"
        value={comment}
        onChangeText={setComment}
      />

      {loading ? (
        <ActivityIndicator size="large" color="#0000ff" />
      ) : (
        <Button title="Kirim Rating" onPress={handleRating} />
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, padding: 20, justifyContent: 'center' },
  title: { fontSize: 24, fontWeight: 'bold', marginBottom: 20, textAlign: 'center' },
  label: { fontSize: 16, marginBottom: 10 },
  scoreContainer: { flexDirection: 'row', justifyContent: 'space-around', marginBottom: 20 },
  input: { borderWidth: 1, borderColor: '#ccc', padding: 10, marginBottom: 20, borderRadius: 5, textAlignVertical: 'top' }
});
