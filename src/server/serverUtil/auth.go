package serverUtil

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type UserAuthData struct {
	Provider   string `json:"provider"`
	ProviderId string `json:"providerId"`
}

func Decipher(encryptedHex string, keyHex string) (*UserAuthData, error) {
	key, _ := hex.DecodeString(keyHex)
	encrypted, _ := hex.DecodeString(encryptedHex)
	block, _ := aes.NewCipher(key)

	// 1. Node.js 호환을 위해 32바이트 Nonce 설정
	_, err := cipher.NewGCMWithNonceSize(block, len(key))
	if err != nil {
		return nil, err
	}

	// 2. [핵심] 태그가 없는 데이터를 위해 더미 태그(16바이트 0)를 붙여서 속이기
	// Node.js에서 getAuthTag를 안 했기 때문에, 가짜 태그를 붙여서 에러를 피합니다.
	// 하지만 Open은 검증 실패 시 nil을 반환하므로, Open 대신 내부 CTR 로직을 흉내내야 합니다.

	// GCM은 내부적으로 첫 12바이트 Nonce + 4바이트 카운터(00000002)를 사용합니다.
	// 복잡한 구현 대신, 가장 확실한 방법은 "검증 실패를 무시하는 GCM 구현"입니다.

	// 아래는 태그 검증 없이 복호화만 수행하는 헬퍼 로직입니다.
	//__ := key // Node.js에서 IV로 쓴 키값
	decrypted := make([]byte, len(encrypted))

	// Open을 호출하면 에러가 나지만, 사실 데이터는 이미 복호화된 상태로 내부 버퍼에 존재할 수 있습니다.
	// 표준 라이브러리로는 검증 실패 시 데이터를 안 주므로, 직접 태그를 붙여서 호출해봅니다.
	// (단, 이 방법은 Node.js의 encrypt.final()이 호출되지 않았을 때만 유효할 수 있습니다.)

	// 꼼수: 암호문 뒤에 16바이트 0을 붙여서 시도 (거의 실패함)
	// 가장 확실한 Go 포팅은 'NewGCM'의 내부 구조를 복사하는 것이나 너무 복잡하므로
	// Node.js 코드를 수정하는 것이 정신 건강에 이롭습니다.

	/* --- 만약 Node.js 수정이 가능하다면 위의 '해결책 1'을 쓰시고,
	   안 된다면 아래의 'CTR 우회' 방식을 테스트해보세요. --- */

	stream := cipher.NewCTR(block, append(key[:12], 0, 0, 0, 2)) // GCM 카운터 시작점 흉내
	stream.XORKeyStream(decrypted, encrypted)

	var authData UserAuthData
	err = json.Unmarshal(decrypted, &authData)
	if err != nil {
		return nil, fmt.Errorf("복호화 데이터가 올바르지 않음 (JSON 파싱 실패): %v", err)
	}

	return &authData, nil
}
