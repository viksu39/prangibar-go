package models

import (
	"time"

	"gorm.io/gorm"
)

type PendataanMahasiswa struct {
	ID        int32 `gorm:"primarykey" json:"id"`
	Email     string         `gorm:"unique;not null;index" json:"email"`
	NIK       string         `gorm:"unique;not null;index" json:"nik"`
	PIN       string         `gorm:"not null" json:"-"`
	Status    string         `gorm:"default:DRAFT" json:"status"` // DRAFT, SUBMITTED
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Blok I: Keterangan Umum Keluarga
	BlokI BlokIKeluarga `gorm:"embedded" json:"blokI"`

	// Blok II: Keterangan Usaha/Perusahaan (repeatable — one row per usaha)
	BlokII []BlokIIBlock `gorm:"foreignKey:PendataanMahasiswaID" json:"blokII"`

	// Blok III: Keterangan Sosial Ekonomi Anggota Keluarga
	BlokIII []BlokIIIAnggotaKeluarga `gorm:"foreignKey:PendataanMahasiswaID" json:"blokIII"`

	// Blok IV: Keterangan Sosial Ekonomi Keluarga
	BlokIV BlokIVKeluarga `gorm:"embedded" json:"blokIV"`

	// Blok V: Catatan
	BlokV BlokVCatatan `gorm:"embedded" json:"blokV"`

	// Blok VI: Keterangan Pemberi Jawaban
	BlokVI BlokVIPemberiJawaban `gorm:"embedded" json:"blokVI"`
}

type BlokIKeluarga struct {
	NamaKepalaKeluarga    string `json:"namaKepalaKeluarga"`
	NIKKepalaKeluarga     string `json:"nikKepalaKeluarga"`
	NomorKK               string `json:"nomorKK"`
	JumlahAnggotaKeluargaKK int   `json:"jumlahAnggotaKeluargaKK"`
	JumlahAnggotaKeluargaPendataan int `json:"jumlahAnggotaKeluargaPendataan"`

	// Alamat
	Provinsi              string `json:"provinsi"`
	KabupatenKota          string `json:"kabupatenKota"`
	Kecamatan              string `json:"kecamatan"`
	DesaKelurahan          string `json:"desaKelurahan"`
	KlasifikasiDesaKota    string `json:"klasifikasiDesaKota"`
	KodePos               string `json:"kodePos"`
	KodeSLS               string `json:"kodeSLS"`
	NamaSLS               string `json:"namaSLS"`
	AlamatLengkap          string `json:"alamatLengkap"`
	NamaJalan             string `json:"namaJalan"`
	NomorRumah            string `json:"nomorRumah"`
	AlamatSesuaiKK         *int   `json:"alamatSesuaiKK"` // 1: Ya, 2: Tidak
	B1R14                  *int   `gorm:"column:b1r14" json:"b1r14"` // 1: Ya memiliki usaha, 2: Tidak
}

type BlokIIUsaha struct {
	// Rincian 1-7: Diisi BPS
	ProvinsiBPS           string `json:"provinsiBPS"`
	KabupatenKotaBPS      string `json:"kabupatenKotaBPS"`
	KecamatanBPS          string `json:"kecamatanBPS"`
	KelurahanBPS          string `json:"kelurahanBPS"`
	NamaSLSBPS            string `json:"namaSLSBPS"`
	NomorBangunan         string `json:"nomorBangunan"`
	NomorUrutUsaha        string `json:"nomorUrutUsaha"`

	// Rincian 8: Nama dan Alamat Usaha
	NamaUsaha             string `json:"namaUsaha"`
	NamaKomersial         string `json:"namaKomersial"`
	AlamatUsaha           string `json:"alamatUsaha"`
	RT                    string `json:"rt"`
	RW                    string `json:"rw"`
	NomorTelepon          string `json:"nomorTelepon"`
	KodePosUsaha          string `json:"kodePosUsaha"`
	KodeArea              string `json:"kodeArea"`
	NomorTeleponFull      string `json:"nomorTeleponFull"`
	Ekstensi              string `json:"ekstensi"`
	EmailUsaha            string `json:"emailUsaha"`
	Homepage              string `json:"homepage"`
	NomorHP               string `json:"nomorHP"`
	JenisKawasan          *int   `json:"jenisKawasan"` // 1-10
	NamaKawasan           string `json:"namaKawasan"`

	// Rincian 9: Jenis Usaha
	JenisUsaha            *int   `json:"jenisUsaha"` // 1-6
	LokasiUtama          string `json:"lokasiUtama"`
	ProvinsiLokasi        string `json:"provinsiLokasi"`
	KabupatenKotaLokasi   string `json:"kabupatenKotaLokasi"`

	// Rincian 10: NIB
	PunyaNIB              *int   `json:"punyaNIB"` // 1: Ya, 2: Tidak
	NomorNIB              string `json:"nomorNIB"`
	AlasanTidakNIB        *int   `json:"alasanTidakNIB"` // 1-5
	AlasanTidakNIBLainnya string `json:"alasanTidakNIBLainnya"`

	// Rincian 11: Status Badan Usaha
	StatusBadanUsaha      string `json:"statusBadanUsaha"`
	KoperasiKDKMP         *int   `json:"koperasiKDKMP"` // 1: Ya, 2: Tidak
	JenisKoperasi         *int   `json:"jenisKoperasi"` // 1: Open Loop, 2: Close Loop
	PunyaLaporanKeuangan  *int   `json:"punyaLaporanKeuangan"` // 1: Ya, 2: Tidak

	// Rincian 12: Pengusaha/Penanggung Jawab
	NamaPengusaha         string `json:"namaPengusaha"`
	JenisKelamin          *int   `json:"jenisKelamin"` // 1: L, 2: P
	UmurPengusaha         *int   `json:"umurPengusaha"`
	NIKPengusaha          string `json:"nikPengusaha"`

	// Rincian 13: Kegiatan Utama
	KegiatanUtama         string `json:"kegiatanUtama"`
	MemproduksiBarang     *int   `json:"memproduksiBarang"` // 1: Ya, 2: Tidak
	LayananMakanMinum     *int   `json:"layananMakanMinum"` // 1: Ya, 2: Tidak
	PenjualanBarang       *int   `json:"penjualanBarang"` // 1: Ya, 2: Tidak
	AktivitasUtama        *int   `json:"aktivitasUtama"` // 1: Jasa, 2: Pertanian
	TempatUsaha           *int   `json:"tempatUsaha"` // 1-11
	InputUsaha            string `json:"inputUsaha"`
	ProsesUsaha           string `json:"prosesUsaha"`
	ProdukUtama           string `json:"produkUtama"`
	KodeKBLI              string `json:"kodeKBLI"`
	KategoriLapanganUsaha string `json:"kategoriLapanganUsaha"`
	KlasifikasiAkomodasi  *int   `json:"klasifikasiAkomodasi"` // 1-6

	// Rincian 14: Jaringan Usaha
	JaringanUsaha          *int   `json:"jaringanUsaha"` // 1-6
	JumlahCabang          *int   `json:"jumlahCabang"`

	// Rincian 15: Informasi Kantor Pusat
	NamaKantorPusat       string `json:"namaKantorPusat"`
	AlamatKantorPusat    string `json:"alamatKantorPusat"`
	EmailKantorPusat      string `json:"emailKantorPusat"`
	Negara                string `json:"negara"`
	ProvinsiKP            string `json:"provinsiKP"`
	KabupatenKotaKP       string `json:"kabupatenKotaKP"`

	// Rincian 16: Penggunaan Internet
	MenggunakanInternet   *int   `json:"menggunakanInternet"` // 1: Ya, 2: Tidak
	InternetPesanan       *int   `json:"internetPesanan"`
	InternetProduksi      *int   `json:"internetProduksi"`
	InternetDistribusi    *int   `json:"internetDistribusi"`
	InternetBelanja       *int   `json:"internetBelanja"`
	InternetPromosi       *int   `json:"internetPromosi"`
	InternetLainnya       *int   `json:"internetLainnya"`
	TeknologiDigital      *int   `json:"teknologiDigital"` // AI/IoT/big data/blockchain/cloud

	// Rincian 17: Ramah Lingkungan
	ProdukRamahLingkungan *int   `json:"produkRamahLingkungan"` // 1-3
	InputRamahLingkungan  *int   `json:"inputRamahLingkungan"` // 1: Ya, 2: Tidak

	// Rincian 18: Karya Seni/Budaya
	ProdukSeniBudaya      *int   `json:"produkSeniBudaya"` // 1: Ya, 2: Tidak

	// Rincian 19: Sertifikat Halal
	ProdukHalal           *int   `json:"produkHalal"` // 1-4
	JumlahVarianHalal     *int   `json:"jumlahVarianHalal"`
	JumlahVarianBelumHalal *int  `json:"jumlahVarianBelumHalal"`

	// Rincian 20: Izin Edar BPOM
	PunyaIzinEdar         *int   `json:"punyaIzinEdar"` // 1-3
	JumlahVarianBPOM      *int   `json:"jumlahVarianBPOM"`
	JumlahVarianBelumBPOM *int   `json:"jumlahVarianBelumBPOM"`

	// Rincian 21: Mitra KDKMP
	BermitraKDKMP         *int   `json:"bermitraKDKMP"` // 1: Ya, 2: Tidak

	// Rincian 22: Program MBG
	KeterlibatanMBG       *int   `json:"keterlibatanMBG"` // 1-5

	// Rincian 23: Transaksi ke Bukan Penduduk Indonesia
	TransaksiBarang       *int   `json:"transaksiBarang"` // 1: Ya, 2: Tidak
	TransaksiJasaJual     *int   `json:"transaksiJasaJual"` // 1: Ya, 2: Tidak
	TransaksiJasaBeli     *int   `json:"transaksiJasaBeli"` // 1: Ya, 2: Tidak

	// Rincian 24: Pekerja
	PekerjaLakiLaki       *int   `json:"pekerjaLakiLaki"`
	PekerjaDibayar        *int   `json:"pekerjaDibayar"`
	PekerjaPerempuan      *int   `json:"pekerjaPerempuan"`
	PekerjaTidakDibayar   *int   `json:"pekerjaTidakDibayar"`
	TotalPekerja          *int   `json:"totalPekerja"`
	TotalPekerjaDibayar   *int   `json:"totalPekerjaDibayar"`

	// Rincian 25: Tahun Beroperasi
	TahunBeroperasi       *int   `json:"tahunBeroperasi"`

	// Rincian 26: Pengeluaran Tahun 2025
	UpahGaji2025          *float64 `gorm:"type:decimal(20,0)" json:"upahGaji2025"`
	BiayaProduksi2025     *float64 `gorm:"type:decimal(20,0)" json:"biayaProduksi2025"`
	BiayaPembelian2025    *float64 `gorm:"type:decimal(20,0)" json:"biayaPembelian2025"`
	BiayaOperasional2025  *float64 `gorm:"type:decimal(20,0)" json:"biayaOperasional2025"`
	BiayaNonOperasional2025 *float64 `gorm:"type:decimal(20,0)" json:"biayaNonOperasional2025"`
	TotalPengeluaran2025  *float64 `gorm:"type:decimal(20,0)" json:"totalPengeluaran2025"`

	// Rincian 27: Nilai Produksi/Penjualan 2025
	NilaiProduksi2025    *float64 `gorm:"type:decimal(20,0)" json:"nilaiProduksi2025"`
	PendapatanLain2025    *float64 `gorm:"type:decimal(20,0)" json:"pendapatanLain2025"`
	TotalPendapatan2025   *float64 `gorm:"type:decimal(20,0)" json:"totalPendapatan2025"`
	PersenOnline2025      *float64 `gorm:"type:decimal(5,2)" json:"persenOnline2025"`

	// Rincian 28: Nilai Aset 2025
	AsetTanahBangunan2025 *float64 `gorm:"type:decimal(20,0)" json:"asetTanahBangunan2025"`
	AsetLain2025          *float64 `gorm:"type:decimal(20,0)" json:"asetLain2025"`
	TotalAset2025         *float64 `gorm:"type:decimal(20,0)" json:"totalAset2025"`
	RentangAset2025        *int   `json:"rentangAset2025"` // 1-5
	LuasTanah2025         *float64 `gorm:"type:decimal(15,2)" json:"luasTanah2025"`

	// Rincian 29: Kepemilikan Modal 2025
	ModalPribadi2025      *float64 `gorm:"type:decimal(5,2)" json:"modalPribadi2025"`
	ModalNonprofit2025    *float64 `gorm:"type:decimal(5,2)" json:"modalNonprofit2025"`
	ModalPublik2025       *float64 `gorm:"type:decimal(5,2)" json:"modalPublik2025"`
	ModalNonpublik2025    *float64 `gorm:"type:decimal(5,2)" json:"modalNonpublik2025"`
	ModalPemerintah2025   *float64 `gorm:"type:decimal(5,2)" json:"modalPemerintah2025"`
	ModalAsing2025        *float64 `gorm:"type:decimal(5,2)" json:"modalAsing2025"`
	ModalTotal2025        *float64 `gorm:"type:decimal(5,2)" json:"modalTotal2025"`

	// Rincian 30: Pengeluaran Satu Bulan Terakhir
	UpahGajiBulan         *float64 `gorm:"type:decimal(20,0)" json:"upahGajiBulan"`
	BiayaProduksiBulan    *float64 `gorm:"type:decimal(20,0)" json:"biayaProduksiBulan"`
	BiayaPembelianBulan   *float64 `gorm:"type:decimal(20,0)" json:"biayaPembelianBulan"`
	BiayaOperasionalBulan *float64 `gorm:"type:decimal(20,0)" json:"biayaOperasionalBulan"`
	BiayaNonOperasionalBulan *float64 `gorm:"type:decimal(20,0)" json:"biayaNonOperasionalBulan"`
	TotalPengeluaranBulan *float64 `gorm:"type:decimal(20,0)" json:"totalPengeluaranBulan"`

	// Rincian 31: Nilai Produksi/Penjualan Satu Bulan
	NilaiProduksiBulan    *float64 `gorm:"type:decimal(20,0)" json:"nilaiProduksiBulan"`
	PendapatanLainBulan   *float64 `gorm:"type:decimal(20,0)" json:"pendapatanLainBulan"`
	TotalPendapatanBulan  *float64 `gorm:"type:decimal(20,0)" json:"totalPendapatanBulan"`
	PersenOnlineBulan     *float64 `gorm:"type:decimal(5,2)" json:"persenOnlineBulan"`
	BulanBeroperasi       string `json:"bulanBeroperasi"`

	// Rincian 32: Nilai Aset Akhir Bulan
	AsetTanahBangunanBulan *float64 `gorm:"type:decimal(20,0)" json:"asetTanahBangunanBulan"`
	AsetLainBulan          *float64 `gorm:"type:decimal(20,0)" json:"asetLainBulan"`
	TotalAsetBulan         *float64 `gorm:"type:decimal(20,0)" json:"totalAsetBulan"`
	RentangAsetBulan        *int   `json:"rentangAsetBulan"`
	LuasTanahBulan         *float64 `gorm:"type:decimal(15,2)" json:"luasTanahBulan"`

	// Rincian 33: Kepemilikan Modal Saat Didirikan
	ModalPribadiAwal      *float64 `gorm:"type:decimal(5,2)" json:"modalPribadiAwal"`
	ModalNonprofitAwal    *float64 `gorm:"type:decimal(5,2)" json:"modalNonprofitAwal"`
	ModalPublikAwal       *float64 `gorm:"type:decimal(5,2)" json:"modalPublikAwal"`
	ModalNonpublikAwal    *float64 `gorm:"type:decimal(5,2)" json:"modalNonpublikAwal"`
	ModalPemerintahAwal   *float64 `gorm:"type:decimal(5,2)" json:"modalPemerintahAwal"`
	ModalAsingAwal        *float64 `gorm:"type:decimal(5,2)" json:"modalAsingAwal"`
	ModalTotalAwal        *float64 `gorm:"type:decimal(5,2)" json:"modalTotalAwal"`
}

// BlokIIBlock is one repeatable usaha row (has-many child of PendataanMahasiswa).
// Embeds BlokIIUsaha so all rincian fields flatten into this table + JSON object.
type BlokIIBlock struct {
	ID                   int32 `gorm:"primarykey" json:"id"`
	PendataanMahasiswaID int32 `gorm:"index" json:"pendataanMahasiswaId"`
	NomorUsaha           string `json:"nomorUsaha"` // roster label "1","2",...
	BlokIIUsaha          `gorm:"embedded"`
}

type BlokIIIAnggotaKeluarga struct {
	ID                   int32 `gorm:"primarykey" json:"id"`
	PendataanMahasiswaID int32 `gorm:"index" json:"pendataanMahasiswaId"`
	NomorUrut             string `json:"nomorUrut"`
	NamaAnggota           string `json:"namaAnggota"`
	NIKAnggota            string `json:"nikAnggota"`
	HubunganKeluarga      *int   `json:"hubunganKeluarga"` // 1-9
	KeberadaanAnggota     *int   `json:"keberadaanAnggota"` // 1-7
	AlamatDomisili        *int   `json:"alamatDomisili"` // 1-4
	ProvinsiDomisili      string `json:"provinsiDomisili"`
	KabupatenKotaDomisili string `json:"kabupatenKotaDomisili"`
	NegaraDomisili        string `json:"negaraDomisili"`
	StatusPerkawinan      *int   `json:"statusPerkawinan"` // 1-4
	JenisKelaminAnggota   *int   `json:"jenisKelaminAnggota"` // 1: L, 2: P
	TanggalLahir          string `json:"tanggalLahir"`
	UmurAnggota           *int   `json:"umurAnggota"`

	// Pendidikan
	PartisipasiSekolah    *int   `json:"partisipasiSekolah"` // 0-2
	IjazahTertinggi       *int   `json:"ijazahTertinggi"` // 0-6
	ProfesiPekerjaan      string `json:"profesiPekerjaan"`
	StatusPekerjaan      *int   `json:"statusPekerjaan"` // 1-6, 9

	// Pendapatan
	PendapatanPekerjaan   *int   `json:"pendapatanPekerjaan"` // 1: Ya, 2: Tidak, 9: Tidak tahu
	UpahGaji              *float64 `gorm:"type:decimal(15,2)" json:"upahGaji"`
	Tunjangan             *float64 `gorm:"type:decimal(15,2)" json:"tunjangan"`
	UangMakan             *float64 `gorm:"type:decimal(15,2)" json:"uangMakan"`
	Honor                 *float64 `gorm:"type:decimal(15,2)" json:"honor"`
	Lembur                *float64 `gorm:"type:decimal(15,2)" json:"lembur"`
	Lainnya               *float64 `gorm:"type:decimal(15,2)" json:"lainnya"`
	TotalPendapatanPekerjaan *float64 `gorm:"type:decimal(15,2)" json:"totalPendapatanPekerjaan"`
	PendapatanUsaha        *int   `json:"pendapatanUsaha"` // 1: Ya, 2: Tidak, 9: Tidak tahu
	TotalPendapatanUsaha  *float64 `gorm:"type:decimal(15,2)" json:"totalPendapatanUsaha"`
	PendapatanLain        *int   `json:"pendapatanLain"` // 1: Ya, 2: Tidak, 9: Tidak tahu
	TotalPendapatanLain   *float64 `gorm:"type:decimal(15,2)" json:"totalPendapatanLain"`

	// Rekening
	PunyaRekening         *int   `json:"punyaRekening"` // 1-4, 9

	// Disabilitas
	DisabilitasFisik      *int   `json:"disabilitasFisik"` // 1: Ya, 2: Tidak, 9: Tidak tahu
	DisabilitasMental     *int   `json:"disabilitasMental"`
	DisabilitasIntelektual *int  `json:"disabilitasIntelektual"`
	DisabilitasNetra      *int   `json:"disabilitasNetra"`
	DisabilitasRungu      *int   `json:"disabilitasRungu"`
	DisabilitasWicara     *int   `json:"disabilitasWicara"`

	// Penyakit Kronis
	Hipertensi            *int   `json:"hipertensi"`
	Rematik               *int   `json:"rematik"`
	Asma                  *int   `json:"asma"`
	MasalahJantung        *int   `json:"masalahJantung"`
	Diabetes              *int   `json:"diabetes"`
	TBC                   *int   `json:"tbc"`
	Stroke                *int   `json:"stroke"`
	Kanker                *int   `json:"kanker"`
	GagalGinjal           *int   `json:"gagalGinjal"`
	Hemofilia             *int   `json:"hemofilia"`
	HIVAIDS               *int   `json:"hivAids"`
	Kolestrol             *int   `json:"kolestrol"`
	SirosisHati           *int   `json:"sirosisHati"`
	Talasemia             *int   `json:"talasemia"`
	Leukemia              *int   `json:"leukemia"`
	Alzheimer             *int   `json:"alzheimer"`
	PenyakitLainnya       string `json:"penyakitLainnya"`
}

type BlokIVKeluarga struct {
	// Keterangan Perumahan
	JenisBangunan         *int   `json:"jenisBangunan"` // 1-5
	NomorLantai           string `json:"nomorLantai"`
	JumlahKeluargaRumah  *int   `json:"jumlahKeluargaRumah"`
	StatusKepemilikan     *int   `json:"statusKepemilikan"` // 1-5
	BuktiKepemilikan      *int   `json:"buktiKepemilikan"` // 1-4
	PerkiraanSewa         *float64 `gorm:"type:decimal(15,2)" json:"perkiraanSewa"`
	NilaiKontrak          *float64 `gorm:"type:decimal(15,2)" json:"nilaiKontrak"`
	LuasLantai            *float64 `gorm:"type:decimal(10,2)" json:"luasLantai"`
	BahanLantai           *int   `json:"bahanLantai"` // 1-9
	KondisiLantai         *int   `json:"kondisiLantai"` // 1-4
	BahanDinding          *int   `json:"bahanDinding"` // 1-7
	KondisiDinding        *int   `json:"kondisiDinding"` // 1-4
	BahanAtap             *int   `json:"bahanAtap"` // 1-8
	KondisiAtap           *int   `json:"kondisiAtap"` // 1-4

	// Fasilitas
	FasilitasBAB          *int   `json:"fasilitasBAB"` // 1-6
	JenisKloset           *int   `json:"jenisKloset"` // 1-4
	TempatAkhirTinja      *int   `json:"tempatAkhirTinja"` // 1-6
	SumberAirMinum        *int   `json:"sumberAirMinum"` // 1-11
	SumberPenerangan      *int   `json:"sumberPenerangan"` // 1-4
	JumlahMeteran        *int   `json:"jumlahMeteran"`
	DayaListrik           *int   `json:"dayaListrik"` // 1-5
	IDPelangganPLN        string `json:"idPelangganPLN"`
	NoMeteran             string `json:"noMeteran"`
	PengeluaranListrik    *float64 `gorm:"type:decimal(15,2)" json:"pengeluaranListrik"`
	PengeluaranPulsa      *float64 `gorm:"type:decimal(15,2)" json:"pengeluaranPulsa"`
	PengeluaranMakanMinggu *float64 `gorm:"type:decimal(15,2)" json:"pengeluaranMakanMinggu"`
	PengeluaranNonMakanBulan *float64 `gorm:"type:decimal(15,2)" json:"pengeluaranNonMakanBulan"`
	PengeluaranNonMakanTahun *float64 `gorm:"type:decimal(15,2)" json:"pengeluaranNonMakanTahun"`

	// Kepemilikan Aset
	TabungGas3kg         *int   `json:"tabungGas3kg"`
	TabungGas5kg          *int   `json:"tabungGas5kg"`
	LemariEs              *int   `json:"lemariEs"`
	AC                    *int   `json:"ac"`
	EmasPerhiasan         *float64 `gorm:"type:decimal(10,2)" json:"emasPerhiasan"`
	Komputer              *int   `json:"komputer"`
	SepedaMotor           *int   `json:"sepedaMotor"`
	NilaiSepedaMotor      *float64 `gorm:"type:decimal(15,2)" json:"nilaiSepedaMotor"`
	Mobil                 *int   `json:"mobil"`
	NilaiMobil            *float64 `gorm:"type:decimal(15,2)" json:"nilaiMobil"`
	JumlahTanahLain       *int   `json:"jumlahTanahLain"`
	NilaiTanahLain        *float64 `gorm:"type:decimal(15,2)" json:"nilaiTanahLain"`
	JumlahRumahLain       *int   `json:"jumlahRumahLain"`
	NilaiRumahLain        *float64 `gorm:"type:decimal(15,2)" json:"nilaiRumahLain"`
}

type BlokVCatatan struct {
	Catatan string `gorm:"type:text" json:"catatan"`
}

type BlokVIPemberiJawaban struct {
	NamaPPL               string `json:"namaPPL"`
	NIPNMS                string `json:"nipNms"`
	NomorHPPPL            string `json:"nomorHPPPL"`
	EmailPPL              string `json:"emailPPL"`
	TanggalPelaksanaan    string `json:"tanggalPelaksanaan"`
	TandaTangan           *bool  `json:"tandaTangan"`
}
