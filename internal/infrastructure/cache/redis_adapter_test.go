package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/rodrigocostarcs/pix-generator/internal/infrastructure/cache"
	"github.com/stretchr/testify/assert"
)

func TestRedisAdapter(t *testing.T) {
	// Inicializar um servidor Redis em memória para testes
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Não foi possível iniciar o miniredis: %v", err)
	}
	defer mr.Close()

	// Criar um adaptador Redis apontando para o miniredis
	adapter := cache.NewRedisAdapter(mr.Host(), mr.Port(), "", 0)
	ctx := context.Background()

	t.Run("Set_Get_String", func(t *testing.T) {
		// Armazenar um valor simples
		key := "test:string"
		value := "valor de teste"
		err := adapter.Set(ctx, key, value, time.Minute)
		assert.NoError(t, err)

		// Recuperar o valor
		result, err := adapter.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, string(result))
	})

	t.Run("Set_Get_Bytes", func(t *testing.T) {
		// Armazenar bytes
		key := "test:bytes"
		value := []byte{1, 2, 3, 4, 5}
		err := adapter.Set(ctx, key, value, time.Minute)
		assert.NoError(t, err)

		// Recuperar o valor
		result, err := adapter.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("Get_NonExistent", func(t *testing.T) {
		// Tentar recuperar uma chave que não existe
		key := "test:nonexistent"
		_, err := adapter.Get(ctx, key)
		assert.Error(t, err)
	})

	t.Run("SetWithExpiration", func(t *testing.T) {
		// Armazenar com expiração curta
		key := "test:expiration"
		value := "valor temporário"
		err := adapter.Set(ctx, key, value, 50*time.Millisecond)
		assert.NoError(t, err)

		// Verificar que o valor existe agora
		result, err := adapter.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, string(result))

		// Esperar a expiração
		time.Sleep(100 * time.Millisecond)
		mr.FastForward(100 * time.Millisecond)

		// Verificar que o valor não existe mais
		_, err = adapter.Get(ctx, key)
		assert.Error(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		// Armazenar um valor
		key := "test:delete"
		value := "valor para deletar"
		err := adapter.Set(ctx, key, value, time.Minute)
		assert.NoError(t, err)

		// Verificar que o valor existe
		result, err := adapter.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, string(result))

		// Deletar o valor
		err = adapter.Delete(ctx, key)
		assert.NoError(t, err)

		// Verificar que o valor não existe mais
		_, err = adapter.Get(ctx, key)
		assert.Error(t, err)
	})

	t.Run("GetObject_SetObject", func(t *testing.T) {
		// Estrutura para teste
		type TestStruct struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Value bool   `json:"value"`
		}

		// Criar objeto
		testObject := TestStruct{
			ID:    123,
			Name:  "Teste Object",
			Value: true,
		}

		key := "test:object"

		// Armazenar objeto
		err := cache.SetObject(adapter, ctx, key, testObject, time.Minute)
		assert.NoError(t, err)

		// Recuperar objeto
		var retrievedObject TestStruct
		err = cache.GetObject(adapter, ctx, key, &retrievedObject)
		assert.NoError(t, err)

		// Verificar valores
		assert.Equal(t, testObject.ID, retrievedObject.ID)
		assert.Equal(t, testObject.Name, retrievedObject.Name)
		assert.Equal(t, testObject.Value, retrievedObject.Value)
	})
}
