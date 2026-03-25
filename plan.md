# 마이그레이션 계획 (Migration Plan)

기존 Node.js 구현(`.old/`)과 새로운 Go 구현(`src/`)을 비교하여 누락된 기능 및 불일치 사항을 정리합니다.

## 1. 인증 및 보안 (Authentication & Security)
- [x] **AES-256-GCM 복호화 구현**: `auth-user` 쿠키를 복호화하기 위한 Go 버전의 `decipher` 함수 구현.
- [x] **인증 확인 로직 추가**:
    - [x] `POST /upload/img`: `auth-user` 쿠키 확인 -> 복호화 -> 유저 등급(grade >= 9) 검증.
    - [x] `POST /upload/link`: `key`(바디의 API Key) 또는 `auth-user` 쿠키 지원. 쿠키 사용 시 유저 등급 검증.
- [x] **API Key 검증**: `link-uploader.go`에서 `requestData.key`와 `API_KEY` 환경 변수 비교 로직 추가.

## 2. 데이터베이스 통합 (Database Integration)
- [x] **유저 데이터 조회 함수 추가**: `src/server/db/db.go`에 `provider`와 `providerId`로 유저를 찾는 `GetUserDataByProvider` 함수 추가.
- [x] **파일 로그 기록**: `img-uploader.go`와 `link-uploader.go`에 `db.NewFileLog` 연동.

## 3. 라우트 및 경로 수정 (Route & Path Corrections)
- [x] **Fumen 라우트 수정**: `src/server/server.go`의 경로 파라미터 수정 완료.
- [x] **파일 저장 및 조회 경로 일관성 유지**:
    - 모든 파일 서비스가 `util.FileDir()` 하위를 참조하도록 수정 완료.

## 4. 환경 변수 및 기타
- [x] 다음 환경 변수들이 `os.Getenv` 등을 통해 올바르게 로드되는지 확인:
    - `AUTH_KEY`: 쿠키 복호화용.
    - `API_KEY`: 외부 API 업로드 인증용.
- [x] `POST /upload/link`의 에러 핸들링 보강.
