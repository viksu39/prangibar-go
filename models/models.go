package models

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Name      string         `gorm:"not null" json:"name"`
	Password  string         `gorm:"not null" json:"-"`
	CreatedAt time.Time      `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deletedAt" json:"-"`
}

func (Admin) TableName() string { return "admin" }

type Provinsi struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Kode          string         `gorm:"unique;not null" json:"kode"`
	Nama          string         `gorm:"not null" json:"nama"`
	KabupatenKota []KabupatenKota `json:"-"`
}

func (Provinsi) TableName() string { return "provinsi" }

type KabupatenKota struct {
	ID         uint        `gorm:"primarykey" json:"id"`
	Kode       string      `gorm:"not null" json:"kode"`
	Nama       string      `gorm:"not null" json:"nama"`
	ProvinsiID uint        `gorm:"not null;column:provinsiId" json:"provinsiId"`
	Provinsi   Provinsi    `json:"-"`
	Kecamatan  []Kecamatan `json:"-"`
}

func (KabupatenKota) TableName() string { return "kabupatenkota" }

type Kecamatan struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	Kode            string         `gorm:"not null" json:"kode"`
	Nama            string         `gorm:"not null" json:"nama"`
	KabupatenKotaID uint           `gorm:"not null;column:kabupatenKotaId" json:"kabupatenKotaId"`
	KabupatenKota   KabupatenKota  `json:"-"`
	Desa            []Desa         `json:"-"`
}

func (Kecamatan) TableName() string { return "kecamatan" }

type Desa struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Kode        string    `gorm:"not null" json:"kode"`
	Nama        string    `gorm:"not null" json:"nama"`
	Klas        *int      `json:"klas"`
	KecamatanID uint      `gorm:"not null;column:kecamatanId" json:"kecamatanId"`
	Kecamatan   Kecamatan `json:"-"`
}

func (Desa) TableName() string { return "desa" }

type StatusPendataan string

const (
	StatusBelum   StatusPendataan = "BELUM"
	StatusSelesai StatusPendataan = "SELESAI"
)

type SkalaUsaha string

const (
	SkalaUB    SkalaUsaha = "UB"
	SkalaUMKM  SkalaUsaha = "UMKM"
)

type Perusahaan struct {
	ID            uint             `gorm:"primarykey" json:"id"`
	Nama          string           `gorm:"not null" json:"nama"`
	Alamat        *string          `json:"alamat"`
	ContactPerson *string          `gorm:"column:contactPerson" json:"contactPerson"`
	Email         *string          `json:"email"`
	Phone         *string          `json:"phone"`
	B1R1          *string          `gorm:"type:varchar(2);column:b1r1" json:"b1r1"`
	B1R2          *string          `gorm:"type:varchar(4);column:b1r2" json:"b1r2"`
	B1R3          *string          `gorm:"type:varchar(7);column:b1r3" json:"b1r3"`
	B1R4          *string          `gorm:"type:varchar(10);column:b1r4" json:"b1r4"`
	Token         *string          `gorm:"unique" json:"token"`
	Status        StatusPendataan  `gorm:"default:BELUM" json:"status"`
	SkalaUsaha    *SkalaUsaha      `gorm:"column:skalaUsaha" json:"skalaUsaha"`
	CreatedAt     time.Time        `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt     time.Time        `gorm:"column:updatedAt" json:"updatedAt"`
	Pendataan     *Pendataan       `json:"pendataan,omitempty"`
	PendataanUmkm *PendataanUmkm   `json:"pendataanUmkm,omitempty"`
}

func (Perusahaan) TableName() string { return "perusahaan" }

type Pendataan struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	PerusahaanID uint      `gorm:"unique;not null;column:perusahaanId" json:"perusahaanId"`
	Perusahaan   Perusahaan `json:"perusahaan,omitempty"`

	B1R1  *string `gorm:"type:varchar(2)" json:"b1r1"`
	B1R2  *string `gorm:"type:varchar(4)" json:"b1r2"`
	B1R3  *string `gorm:"type:varchar(7)" json:"b1r3"`
	B1R4  *string `gorm:"type:varchar(10)" json:"b1r4"`

	B1R5A            *string `json:"b1r5a"`
	B1R5B            *string `json:"b1r5b"`
	B1R5CAlamat      *string `json:"b1r5cAlamat"`
	B1R5CRT          *string `json:"b1r5cRt"`
	B1R5CRW          *string `json:"b1r5cRw"`
	B1R5CKodePos     *string `json:"b1r5cKodePos"`
	B1R5CTelKodeArea *string `json:"b1r5cTelKodeArea"`
	B1R5CTelNomor    *string `json:"b1r5cTelNomor"`
	B1R5CTelEkstensi *string `json:"b1r5cTelEkstensi"`
	B1R5CEmail       *string `json:"b1r5cEmail"`
	B1R5CHomepage    *string `json:"b1r5cHomepage"`
	B1R5CHP          *string `json:"b1r5cHp"`
	B1R5D            *int    `json:"b1r5d"`
	B1R5E            *string `json:"b1r5e"`

	B1R6A        *int    `json:"b1r6a"`
	B1R6B        *string `json:"b1r6b"`
	B1R6C        *int    `json:"b1r6c"`
	B1R6CLainnya *string `json:"b1r6cLainnya"`

	B1R7A *int `json:"b1r7a"`
	B1R7B *int `json:"b1r7b"`
	B1R7C *int `json:"b1r7c"`
	B1R7D *int `json:"b1r7d"`

	B1R8A *string `json:"b1r8a"`
	B1R8B *int    `json:"b1r8b"`
	B1R8C *int    `json:"b1r8c"`
	B1R8D *string `json:"b1r8d"`

	B1R9A  *string `json:"b1r9a"`
	B1R9B1 *int    `json:"b1r9b1"`
	B1R9B2 *int    `json:"b1r9b2"`
	B1R9B3 *int    `json:"b1r9b3"`
	B1R9B4 *int    `json:"b1r9b4"`
	B1R9C  *int    `json:"b1r9c"`
	B1R9D  *string `json:"b1r9d"`
	B1R9E  *string `json:"b1r9e"`
	B1R9F  *string `json:"b1r9f"`
	B1R9G  *string `json:"b1r9g"`
	B1R9H  *string `json:"b1r9h"`
	B1R9I  *int    `json:"b1r9i"`

	B1R10A *int `json:"b1r10a"`
	B1R10B *int `json:"b1r10b"`

	B1R11A *string `json:"b1r11a"`
	B1R11B *string `json:"b1r11b"`
	B1R11C *string `json:"b1r11c"`
	B1R11D *string `json:"b1r11d"`
	B1R11E *string `json:"b1r11e"`
	B1R11F *string `json:"b1r11f"`

	B1R12A  *int `json:"b1r12a"`
	B1R12B1 *int `json:"b1r12b1"`
	B1R12B2 *int `json:"b1r12b2"`
	B1R12B3 *int `json:"b1r12b3"`
	B1R12B4 *int `json:"b1r12b4"`
	B1R12B5 *int `json:"b1r12b5"`
	B1R12B6 *int `json:"b1r12b6"`
	B1R12C  *int `json:"b1r12c"`

	B1R13A *int `json:"b1r13a"`
	B1R13B *int `json:"b1r13b"`

	B1R14 *int `json:"b1r14"`

	B1R15A *int `json:"b1r15a"`
	B1R15B *int `json:"b1r15b"`
	B1R15C *int `json:"b1r15c"`

	B1R16A *int `json:"b1r16a"`
	B1R16B *int `json:"b1r16b"`
	B1R16C *int `json:"b1r16c"`

	B1R17 *int `json:"b1r17"`

	B1R18 *int `json:"b1r18"`

	B1R19A *int `json:"b1r19a"`
	B1R19B *int `json:"b1r19b"`

	B1R20A *int `json:"b1r20a"`
	B1R20B *int `json:"b1r20b"`
	B1R20C *int `json:"b1r20c"`

	B1R21 *int `json:"b1r21"`

	B1R22A *float64 `gorm:"type:decimal(20,0)" json:"b1r22a"`
	B1R22B *float64 `gorm:"type:decimal(20,0)" json:"b1r22b"`
	B1R22C *float64 `gorm:"type:decimal(20,0)" json:"b1r22c"`
	B1R22D *float64 `gorm:"type:decimal(20,0)" json:"b1r22d"`
	B1R22E *float64 `gorm:"type:decimal(20,0)" json:"b1r22e"`
	B1R22F *float64 `gorm:"type:decimal(20,0)" json:"b1r22f"`

	B1R23A *float64 `gorm:"type:decimal(20,0)" json:"b1r23a"`
	B1R23B *float64 `gorm:"type:decimal(20,0)" json:"b1r23b"`
	B1R23C *float64 `gorm:"type:decimal(20,0)" json:"b1r23c"`
	B1R23D *float64 `gorm:"type:decimal(5,2)" json:"b1r23d"`

	B1R24A  *float64 `gorm:"type:decimal(20,0)" json:"b1r24a"`
	B1R24B  *float64 `gorm:"type:decimal(20,0)" json:"b1r24b"`
	B1R24C  *float64 `gorm:"type:decimal(20,0)" json:"b1r24c"`
	B1R24C1 *int     `json:"b1r24c1"`
	B1R24D  *float64 `gorm:"type:decimal(15,2)" json:"b1r24d"`

	B1R25A *float64 `gorm:"type:decimal(5,2)" json:"b1r25a"`
	B1R25B *float64 `gorm:"type:decimal(5,2)" json:"b1r25b"`
	B1R25C *float64 `gorm:"type:decimal(5,2)" json:"b1r25c"`
	B1R25D *float64 `gorm:"type:decimal(5,2)" json:"b1r25d"`
	B1R25E *float64 `gorm:"type:decimal(5,2)" json:"b1r25e"`
	B1R25F *float64 `gorm:"type:decimal(5,2)" json:"b1r25f"`
	B1R25G *float64 `gorm:"type:decimal(5,2)" json:"b1r25g"`

	B2Catatan *string `gorm:"type:text" json:"b2catatan"`

	B3R1 *string `gorm:"type:varchar(100)" json:"b3r1"`
	B3R3 *string `gorm:"type:varchar(20)" json:"b3r3"`
	B3R4 *string `gorm:"type:varchar(100)" json:"b3r4"`

	CatatanNgibar      *string    `gorm:"type:text" json:"catatanNgibar"`
	TanggalPelaksanaan *time.Time `json:"tanggalPelaksanaan"`
	TTD                *bool      `json:"ttd"`

	SubmittedAt time.Time `json:"submittedAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (Pendataan) TableName() string { return "pendataan" }

type PendataanUmkm struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	PerusahaanID uint      `gorm:"unique;not null" json:"perusahaanId"`
	Perusahaan   Perusahaan `json:"perusahaan,omitempty"`

	B1R1 *string `gorm:"type:varchar(2)" json:"b1r1"`
	B1R2 *string `gorm:"type:varchar(4)" json:"b1r2"`
	B1R3 *string `gorm:"type:varchar(7)" json:"b1r3"`
	B1R4 *string `gorm:"type:varchar(10)" json:"b1r4"`

	B1R5 *string `json:"b1r5"`
	B1R6 *string `json:"b1r6"`
	B1R7 *string `json:"b1r7"`

	B1R8A            *string `json:"b1r8a"`
	B1R8B            *string `json:"b1r8b"`
	B1R8CAlamat      *string `json:"b1r8cAlamat"`
	B1R8CRT          *string `json:"b1r8cRt"`
	B1R8CRW          *string `json:"b1r8cRw"`
	B1R8CKodePos     *string `json:"b1r8cKodePos"`
	B1R8CTelKodeArea *string `json:"b1r8cTelKodeArea"`
	B1R8CTelNomor    *string `json:"b1r8cTelNomor"`
	B1R8CTelEkstensi *string `json:"b1r8cTelEkstensi"`
	B1R8CEmail       *string `json:"b1r8cEmail"`
	B1R8CHomepage    *string `json:"b1r8cHomepage"`
	B1R8CHP          *string `json:"b1r8cHp"`
	B1R8D            *int    `json:"b1r8d"`
	B1R8E            *string `json:"b1r8e"`

	B1R9A  *int    `json:"b1r9a"`
	B1R9B1 *string `json:"b1r9b1"`
	B1R9B2 *string `json:"b1r9b2"`
	B1R9B3 *string `json:"b1r9b3"`

	B1R10A        *int    `json:"b1r10a"`
	B1R10B        *string `json:"b1r10b"`
	B1R10C        *int    `json:"b1r10c"`
	B1R10CLainnya *string `json:"b1r10cLainnya"`

	B1R11A *string `json:"b1r11a"`
	B1R11B *int    `json:"b1r11b"`
	B1R11C *int    `json:"b1r11c"`
	B1R11D *int    `json:"b1r11d"`

	B1R12A *string `json:"b1r12a"`
	B1R12B *int    `json:"b1r12b"`
	B1R12C *int    `json:"b1r12c"`
	B1R12D *string `gorm:"type:varchar(16)" json:"b1r12d"`

	B1R13A  *string `json:"b1r13a"`
	B1R13B1 *int    `json:"b1r13b1"`
	B1R13B2 *int    `json:"b1r13b2"`
	B1R13B3 *int    `json:"b1r13b3"`
	B1R13B4 *int    `json:"b1r13b4"`
	B1R13C  *int    `json:"b1r13c"`
	B1R13D  *string `json:"b1r13d"`
	B1R13E  *string `json:"b1r13e"`
	B1R13F  *string `json:"b1r13f"`
	B1R13G  *string `gorm:"type:varchar(5)" json:"b1r13g"`
	B1R13H  *string `gorm:"type:varchar(2)" json:"b1r13h"`
	B1R13I  *int    `json:"b1r13i"`

	B1R14A *int `json:"b1r14a"`
	B1R14B *int `json:"b1r14b"`

	B1R15A *string `json:"b1r15a"`
	B1R15B *string `json:"b1r15b"`
	B1R15C *string `json:"b1r15c"`
	B1R15D *string `json:"b1r15d"`
	B1R15E *string `json:"b1r15e"`
	B1R15F *string `json:"b1r15f"`

	B1R16A  *int `json:"b1r16a"`
	B1R16B1 *int `json:"b1r16b1"`
	B1R16B2 *int `json:"b1r16b2"`
	B1R16B3 *int `json:"b1r16b3"`
	B1R16B4 *int `json:"b1r16b4"`
	B1R16B5 *int `json:"b1r16b5"`
	B1R16B6 *int `json:"b1r16b6"`
	B1R16C  *int `json:"b1r16c"`

	B1R17A *int `json:"b1r17a"`
	B1R17B *int `json:"b1r17b"`

	B1R18 *int `json:"b1r18"`

	B1R19A *int `json:"b1r19a"`
	B1R19B *int `json:"b1r19b"`
	B1R19C *int `json:"b1r19c"`

	B1R20A *int `json:"b1r20a"`
	B1R20B *int `json:"b1r20b"`
	B1R20C *int `json:"b1r20c"`

	B1R21 *int `json:"b1r21"`

	B1R22 *int `json:"b1r22"`

	B1R23A *int `json:"b1r23a"`
	B1R23B *int `json:"b1r23b"`
	B1R23C *int `json:"b1r23c"`

	B1R24A1 *int `json:"b1r24a1"`
	B1R24A2 *int `json:"b1r24a2"`
	B1R24B1 *int `json:"b1r24b1"`
	B1R24B2 *int `json:"b1r24b2"`
	B1R24C1 *int `json:"b1r24c1"`
	B1R24C2 *int `json:"b1r24c2"`

	B1R25 *int `json:"b1r25"`

	B1R26A *float64 `gorm:"type:decimal(20,0)" json:"b1r26a"`
	B1R26B *float64 `gorm:"type:decimal(20,0)" json:"b1r26b"`
	B1R26C *float64 `gorm:"type:decimal(20,0)" json:"b1r26c"`
	B1R26D *float64 `gorm:"type:decimal(20,0)" json:"b1r26d"`
	B1R26E *float64 `gorm:"type:decimal(20,0)" json:"b1r26e"`
	B1R26F *float64 `gorm:"type:decimal(20,0)" json:"b1r26f"`

	B1R27A *float64 `gorm:"type:decimal(20,0)" json:"b1r27a"`
	B1R27B *float64 `gorm:"type:decimal(20,0)" json:"b1r27b"`
	B1R27C *float64 `gorm:"type:decimal(20,0)" json:"b1r27c"`
	B1R27D *float64 `gorm:"type:decimal(5,2)" json:"b1r27d"`

	B1R28A  *float64 `gorm:"type:decimal(20,0)" json:"b1r28a"`
	B1R28B  *float64 `gorm:"type:decimal(20,0)" json:"b1r28b"`
	B1R28C  *float64 `gorm:"type:decimal(20,0)" json:"b1r28c"`
	B1R28C1 *int     `json:"b1r28c1"`
	B1R28D  *float64 `gorm:"type:decimal(15,2)" json:"b1r28d"`

	B1R29A *float64 `gorm:"type:decimal(5,2)" json:"b1r29a"`
	B1R29B *float64 `gorm:"type:decimal(5,2)" json:"b1r29b"`
	B1R29C *float64 `gorm:"type:decimal(5,2)" json:"b1r29c"`
	B1R29D *float64 `gorm:"type:decimal(5,2)" json:"b1r29d"`
	B1R29E *float64 `gorm:"type:decimal(5,2)" json:"b1r29e"`
	B1R29F *float64 `gorm:"type:decimal(5,2)" json:"b1r29f"`
	B1R29G *float64 `gorm:"type:decimal(5,2)" json:"b1r29g"`

	B1R30A *float64 `gorm:"type:decimal(20,0)" json:"b1r30a"`
	B1R30B *float64 `gorm:"type:decimal(20,0)" json:"b1r30b"`
	B1R30C *float64 `gorm:"type:decimal(20,0)" json:"b1r30c"`
	B1R30D *float64 `gorm:"type:decimal(20,0)" json:"b1r30d"`
	B1R30E *float64 `gorm:"type:decimal(20,0)" json:"b1r30e"`
	B1R30F *float64 `gorm:"type:decimal(20,0)" json:"b1r30f"`

	B1R31A *float64 `gorm:"type:decimal(20,0)" json:"b1r31a"`
	B1R31B *float64 `gorm:"type:decimal(20,0)" json:"b1r31b"`
	B1R31C *float64 `gorm:"type:decimal(20,0)" json:"b1r31c"`
	B1R31D *float64 `gorm:"type:decimal(5,2)" json:"b1r31d"`
	B1R31E *string  `json:"b1r31e"`

	B1R32A  *float64 `gorm:"type:decimal(20,0)" json:"b1r32a"`
	B1R32B  *float64 `gorm:"type:decimal(20,0)" json:"b1r32b"`
	B1R32C  *float64 `gorm:"type:decimal(20,0)" json:"b1r32c"`
	B1R32C1 *int     `json:"b1r32c1"`
	B1R32D  *float64 `gorm:"type:decimal(15,2)" json:"b1r32d"`

	B1R33A *float64 `gorm:"type:decimal(5,2)" json:"b1r33a"`
	B1R33B *float64 `gorm:"type:decimal(5,2)" json:"b1r33b"`
	B1R33C *float64 `gorm:"type:decimal(5,2)" json:"b1r33c"`
	B1R33D *float64 `gorm:"type:decimal(5,2)" json:"b1r33d"`
	B1R33E *float64 `gorm:"type:decimal(5,2)" json:"b1r33e"`
	B1R33F *float64 `gorm:"type:decimal(5,2)" json:"b1r33f"`
	B1R33G *float64 `gorm:"type:decimal(5,2)" json:"b1r33g"`

	B2Catatan *string `gorm:"type:text" json:"b2catatan"`

	B3R1 *string `gorm:"type:varchar(100)" json:"b3r1"`
	B3R3 *string `gorm:"type:varchar(20)" json:"b3r3"`
	B3R4 *string `gorm:"type:varchar(100)" json:"b3r4"`

	CatatanNgibar      *string    `gorm:"type:text" json:"catatanNgibar"`
	TanggalPelaksanaan *time.Time `json:"tanggalPelaksanaan"`
	TTD                *bool      `json:"ttd"`

	SubmittedAt time.Time `json:"submittedAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (PendataanUmkm) TableName() string { return "pendataanumkm" }

type UpdatePrelistRequest struct {
	Nama          *string            `json:"nama"`
	Alamat        *string            `json:"alamat"`
	ContactPerson *string            `json:"contactPerson"`
	Email         *string            `json:"email"`
	Phone         *string            `json:"phone"`
	B1R1          *string            `json:"b1r1"`
	B1R2          *string            `json:"b1r2"`
	B1R3          *string            `json:"b1r3"`
	B1R4          *string            `json:"b1r4"`
	SkalaUsaha    *SkalaUsaha        `json:"skalaUsaha"`
}

type ApiLog struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	Method         string    `gorm:"type:varchar(10);not null" json:"method"`
	Path           string    `gorm:"type:varchar(500);not null" json:"path"`
	StatusCode     int       `gorm:"not null" json:"statusCode"`
	IP             *string   `gorm:"type:varchar(100)" json:"ip"`
	UserAgent      *string   `gorm:"type:varchar(500)" json:"userAgent"`
	AdminID        *uint     `json:"adminId"`
	ResponseTimeMs int       `gorm:"not null" json:"responseTimeMs"`
	CreatedAt      time.Time `gorm:"index" json:"createdAt"`
}

func (ApiLog) TableName() string { return "apilog" }
