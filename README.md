参孝: https://jvns.ca/blog/2022/02/01/a-dns-resolver-in-80-lines-of-go/

## 実装したい
- [x] 複数Aレコードを返す
- [x] CNAMEレコード対応
- [ ] サブドメイン対応
- [ ] 適切なエラーハンドリング
  - Aレコードがなかった場合
- [ ] miekg/dnsを使わずにDNSパケットのパース
