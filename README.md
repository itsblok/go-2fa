# go-2fa

A minimal Go implementation of TOTP (Time-based One-Time Password) with a custom-built QR code generator.

This project demonstrates a full authentication pipeline from secret generation to QR code rendering without external QR libraries.

## Features

- TOTP secret generation (Base32 encoded)
- HMAC-SHA1 based OTP generation (RFC 6238 compatible)
- QR Code generation from scratch (no external QR libraries)
- Finder pattern rendering
- Timing pattern generation
- Data encoding (byte mode)
- Simplified Reed-Solomon ECC implementation
- 8-mask selection with penalty scoring system
- QR terminal renderer

## Architecture

- `totp/totp.go` → TOTP implementation (RFC 6238)
- `totp/qrcode.go` → QR matrix generation and orchestration
- `totp/encoder.go` → byte encoding + ECC pipeline
- `totp/mask.go` → QR mask application
- `totp/patterns.go` → finder & timing patterns
- `totp/rs.go` → Reed-Solomon ECC (GF256)
- `totp/score.go` → mask penalty scoring

## Notes

- QR implementation is simplified version of QR Code Specification
- Version 1 only (21x21 matrix)
- ECC and format information are simplified but structurally aligned with QR spec
- No external QR libraries are used