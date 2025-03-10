package utils

import "strings"

func MaskEmail(email string) string {
	if email == "" {
		return ""
	}

	// Pisahkan email menjadi bagian sebelum dan sesudah '@'
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email // Kembali ke email asli jika format salah
	}

	username := parts[0]
	domain := parts[1]

	// Mask username: tampilkan 3 karakter awal, sisanya jadi '*'
	maskedUsername := ""
	if len(username) <= 3 {
		maskedUsername = username
	} else {
		maskedUsername = username[:3] + strings.Repeat("*", len(username)-3)
	}

	// Mask domain: sembunyikan semua kecuali bagian setelah titik terakhir
	domainParts := strings.Split(domain, ".")
	if len(domainParts) < 2 {
		return maskedUsername + "@" + domain
	}
	maskedDomain := strings.Repeat("*", len(domain)-len(domainParts[len(domainParts)-1])-1) + "." + domainParts[len(domainParts)-1]

	return maskedUsername + "@" + maskedDomain
}
