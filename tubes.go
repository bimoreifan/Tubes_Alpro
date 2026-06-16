package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Karyawan struct {
	id     int
	nama   string
	divisi string
}

type Absensi struct {
	karyawan  Karyawan
	jamMasuk  int
	jamKeluar int
	status    string
}

const NMAX int = 100
const JAM_MASUK_NORMAL int = 800
const JAM_KELUAR_NORMAL int = 1700

type arrAbsensi [NMAX]Absensi

var scanner = bufio.NewScanner(os.Stdin)

// Helper Input

func readLine() string {
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// Baca integer dengan validasi — error jika bukan angka
func readInt(prompt string) int {
	for {
		fmt.Print(prompt)
		raw := readLine()
		val, err := strconv.Atoi(raw)
		if err != nil {
			fmt.Println("  ✗ Input tidak valid! Harus berupa angka bulat (integer). Coba lagi.")
			continue
		}
		return val
	}
}

// Baca string dengan validasi — error jika berupa angka murni
func readString(prompt string) string {
	for {
		fmt.Print(prompt)
		raw := readLine()
		if raw == "" {
			fmt.Println("  ✗ Input tidak boleh kosong. Coba lagi.")
			continue
		}
		_, err := strconv.Atoi(raw)
		if err == nil {
			fmt.Println("  ✗ Input tidak valid! Tidak boleh berupa angka. Harus teks. Coba lagi.")
			continue
		}
		return raw
	}
}

// Baca jam format HHMM dengan validasi range
func readJam(prompt string) int {
	for {
		fmt.Print(prompt)
		raw := readLine()
		if len(raw) != 4 {
			fmt.Println("  ✗ Format salah! Gunakan 4 digit, contoh: 0800 atau 1700. Coba lagi.")
			continue
		}
		hh, err1 := strconv.Atoi(raw[:2])
		mm, err2 := strconv.Atoi(raw[2:])
		if err1 != nil || err2 != nil {
			fmt.Println("  ✗ Input harus berupa angka 4 digit (HHMM). Coba lagi.")
			continue
		}
		if hh < 0 || hh > 23 {
			fmt.Println("  ✗ Jam tidak valid (00-23). Coba lagi.")
			continue
		}
		if mm < 0 || mm > 59 {
			fmt.Println("  ✗ Menit tidak valid (00-59). Coba lagi.")
			continue
		}
		return hh*100 + mm
	}
}

// Cek duplikat ID

func idSudahAda(a arrAbsensi, n int, id int) bool {
	for i := 0; i < n; i++ {
		if a[i].karyawan.id == id {
			return true
		}
	}
	return false
}

// Logika Status & Format

func tentukanStatus(jamMasuk int, hadir bool) string {
	if !hadir {
		return "Tidak Hadir"
	}
	if jamMasuk <= JAM_MASUK_NORMAL {
		return "Hadir"
	}
	return "Terlambat"
}

func formatJam(jam int) string {
	hh := jam / 100
	mm := jam % 100
	return fmt.Sprintf("%02d:%02d", hh, mm)
}

// Rekursi

func totalTerlambat(a arrAbsensi, n int, i int) int {
	if i >= n {
		return 0
	}
	tambah := 0
	if a[i].status == "Terlambat" {
		tambah = 1
	}
	return tambah + totalTerlambat(a, n, i+1)
}

func totalHadir(a arrAbsensi, n int, i int) int {
	if i >= n {
		return 0
	}
	tambah := 0
	if a[i].status == "Hadir" || a[i].status == "Terlambat" {
		tambah = 1
	}
	return tambah + totalHadir(a, n, i+1)
}

// Input

func inputAbsensi(a *arrAbsensi, n *int) {
	fmt.Println("\n╔══════════════════════════════════════╗")
	fmt.Println("║        INPUT DATA ABSENSI            ║")
	fmt.Println("╚══════════════════════════════════════╝")

	jumlah := readInt("Masukkan jumlah data absensi: ")
	if jumlah <= 0 || jumlah > NMAX {
		fmt.Printf("  ✗ Jumlah harus antara 1 sampai %d.\n", NMAX)
		return
	}

	berhasil := 0
	for i := 0; i < jumlah; i++ {
		fmt.Printf("\n--- Data ke-%d ---\n", i+1)

		var id int
		for {
			id = readInt("ID Karyawan  : ")
			if idSudahAda(*a, *n, id) {
				fmt.Printf("  ✗ ID %d sudah terdaftar! Gunakan ID lain.\n", id)
			} else {
				break
			}
		}

		nama := readString("Nama         : ")
		divisi := readString("Divisi       : ")

		var hadirInput string
		for {
			fmt.Print("Hadir? (y/n) : ")
			hadirInput = strings.ToLower(readLine())
			if hadirInput == "y" || hadirInput == "n" {
				break
			}
			fmt.Println("  ✗ Masukkan 'y' untuk ya atau 'n' untuk tidak.")
		}

		idx := *n + berhasil
		a[idx].karyawan.id = id
		a[idx].karyawan.nama = nama
		a[idx].karyawan.divisi = divisi

		if hadirInput == "y" {
			jamMasuk := readJam("Jam Masuk  (HHMM, ex: 0800) : ")
			jamKeluar := readJam("Jam Keluar (HHMM, ex: 1700) : ")
			for jamKeluar <= jamMasuk {
				fmt.Println("  ✗ Jam keluar harus lebih besar dari jam masuk.")
				jamKeluar = readJam("Jam Keluar (HHMM, ex: 1700) : ")
			}
			a[idx].jamMasuk = jamMasuk
			a[idx].jamKeluar = jamKeluar
			a[idx].status = tentukanStatus(jamMasuk, true)
		} else {
			a[idx].jamMasuk = 0
			a[idx].jamKeluar = 0
			a[idx].status = "Tidak Hadir"
		}
		berhasil++
	}
	*n += berhasil
}

// Cetak

func cetakAbsensi(a arrAbsensi, n int) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                     DATA ABSENSI KARYAWAN                               ║")
	fmt.Println("╠══════╦══════════════════╦════════════════╦═══════════════╦═════════════╣")
	fmt.Printf("║ %-4s ║ %-16s ║ %-14s ║ %-13s ║ %-11s ║\n", "ID", "Nama", "Divisi", "Jam", "Status")
	fmt.Println("╠══════╬══════════════════╬════════════════╬═══════════════╬═════════════╣")
	for i := 0; i < n; i++ {
		jamInfo := "-"
		if a[i].status != "Tidak Hadir" {
			jamInfo = formatJam(a[i].jamMasuk) + "-" + formatJam(a[i].jamKeluar)
		}
		fmt.Printf("║ %-4d ║ %-16s ║ %-14s ║ %-13s ║ %-11s ║\n",
			a[i].karyawan.id,
			a[i].karyawan.nama,
			a[i].karyawan.divisi,
			jamInfo,
			a[i].status)
	}
	fmt.Println("╚══════╩══════════════════╩════════════════╩═══════════════╩═════════════╝")
}

func cetakStatistik(a arrAbsensi, n int) {
	hadir := totalHadir(a, n, 0)
	terlambat := totalTerlambat(a, n, 0)
	tidakHadir := n - hadir

	fmt.Println("\n╔══════════════════════════════════╗")
	fmt.Println("║           STATISTIK              ║")
	fmt.Println("╠══════════════════════════════════╣")
	fmt.Printf("║ Total Karyawan  : %-14d ║\n", n)
	fmt.Printf("║ Hadir (tepat)   : %-14d ║\n", hadir-terlambat)
	fmt.Printf("║ Terlambat       : %-14d ║\n", terlambat)
	fmt.Printf("║ Tidak Hadir     : %-14d ║\n", tidakHadir)
	fmt.Println("╚══════════════════════════════════╝")
}

// Min/Max Jam Masuk

func cariEkstrimJam(a arrAbsensi, n int) {
	idxMin := -1
	idxMax := -1

	for i := 0; i < n; i++ {
		if a[i].status == "Tidak Hadir" {
			continue
		}
		if idxMin == -1 || a[i].jamMasuk < a[idxMin].jamMasuk {
			idxMin = i
		}
		if idxMax == -1 || a[i].jamMasuk > a[idxMax].jamMasuk {
			idxMax = i
		}
	}

	fmt.Println("\n╔══════════════════════════════════════════╗")
	fmt.Println("║         NILAI EKSTRIM JAM MASUK          ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	if idxMin == -1 {
		fmt.Println("Tidak ada karyawan yang hadir.")
		return
	}
	fmt.Printf("Datang Paling Awal  : %-16s (ID: %3d) - Jam %s\n",
		a[idxMin].karyawan.nama, a[idxMin].karyawan.id, formatJam(a[idxMin].jamMasuk))
	fmt.Printf("Datang Paling Akhir : %-16s (ID: %3d) - Jam %s\n",
		a[idxMax].karyawan.nama, a[idxMax].karyawan.id, formatJam(a[idxMax].jamMasuk))
}

// Search

func sequentialSearchNama(a arrAbsensi, n int, keyword string) {
	fmt.Printf("\n[Sequential Search] Mencari nama: %s\n", keyword)
	ketemu := false
	for i := 0; i < n; i++ {
		if strings.EqualFold(a[i].karyawan.nama, keyword) {
			if !ketemu {
				fmt.Println("Hasil ditemukan:")
			}
			fmt.Printf("  → ID: %d | Nama: %s | Divisi: %s | Status: %s\n",
				a[i].karyawan.id, a[i].karyawan.nama, a[i].karyawan.divisi, a[i].status)
			ketemu = true
		}
	}
	if !ketemu {
		fmt.Println("  Data tidak ditemukan.")
	}
}

func binarySearchID(a arrAbsensi, n int, target int) int {
	left, right := 0, n-1
	for left <= right {
		mid := left + (right-left)/2
		if a[mid].karyawan.id == target {
			return mid
		} else if a[mid].karyawan.id < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}

func prosedurBinarySearch(a arrAbsensi, n int) {
	target := readInt("\n[Binary Search] Masukkan ID karyawan yang dicari: ")

	sortedArr := a
	selectionSortAsc(&sortedArr, n)

	idx := binarySearchID(sortedArr, n, target)
	if idx != -1 {
		fmt.Println("Data ditemukan:")
		fmt.Printf("  → ID: %d | Nama: %s | Divisi: %s | Status: %s\n",
			sortedArr[idx].karyawan.id,
			sortedArr[idx].karyawan.nama,
			sortedArr[idx].karyawan.divisi,
			sortedArr[idx].status)
	} else {
		fmt.Println("  Data dengan ID tersebut tidak ditemukan.")
	}
}

// Selection Sort (Ascending & Descending by ID)

func selectionSortAsc(a *arrAbsensi, n int) {
	for pass := 0; pass < n-1; pass++ {
		idxMin := pass
		for i := pass + 1; i < n; i++ {
			if a[i].karyawan.id < a[idxMin].karyawan.id {
				idxMin = i
			}
		}
		a[pass], a[idxMin] = a[idxMin], a[pass]
	}
}

func selectionSortDesc(a *arrAbsensi, n int) {
	for pass := 0; pass < n-1; pass++ {
		idxMax := pass
		for i := pass + 1; i < n; i++ {
			if a[i].karyawan.id > a[idxMax].karyawan.id {
				idxMax = i
			}
		}
		a[pass], a[idxMax] = a[idxMax], a[pass]
	}
}

// Insertion Sort (Ascending & Descending by Jam Masuk)

func insertionSortAsc(a *arrAbsensi, n int) {
	for i := 1; i < n; i++ {
		temp := a[i]
		j := i - 1
		for j >= 0 && a[j].jamMasuk > temp.jamMasuk {
			a[j+1] = a[j]
			j--
		}
		a[j+1] = temp
	}
}

func insertionSortDesc(a *arrAbsensi, n int) {
	for i := 1; i < n; i++ {
		temp := a[i]
		j := i - 1
		for j >= 0 && a[j].jamMasuk < temp.jamMasuk {
			a[j+1] = a[j]
			j--
		}
		a[j+1] = temp
	}
}

// Menu & Main

func tampilkanMenu() {
	fmt.Println("\n╔══════════════════════════════════════════════╗")
	fmt.Println("║         SISTEM ABSENSI KARYAWAN              ║")
	fmt.Println("╠══════════════════════════════════════════════╣")
	fmt.Println("║  1. Input Data Absensi                       ║")
	fmt.Println("║  2. Tampilkan Semua Data                     ║")
	fmt.Println("║  3. Statistik Kehadiran (Rekursif)           ║")
	fmt.Println("║  4. Nilai Ekstrim Jam Masuk (Min/Max)        ║")
	fmt.Println("║  5. Cari Karyawan (Sequential Search)        ║")
	fmt.Println("║  6. Cari Karyawan by ID (Binary Search)      ║")
	fmt.Println("║  7. Selection Sort by ID – Ascending         ║")
	fmt.Println("║  8. Selection Sort by ID – Descending        ║")
	fmt.Println("║  9. Insertion Sort by Jam Masuk – Ascending  ║")
	fmt.Println("║ 10. Insertion Sort by Jam Masuk – Descending ║")
	fmt.Println("║  0. Keluar                                   ║")
	fmt.Println("╚══════════════════════════════════════════════╝")
	fmt.Print("Pilihan: ")
}

func main() {
	var a arrAbsensi
	var n int

	fmt.Println("╔══════════════════════════════════════════════╗")
	fmt.Println("║      TUBES ALPRO - SISTEM ABSENSI            ║")
	fmt.Println("║   Mencakup semua materi modul praktikum      ║")
	fmt.Println("╚══════════════════════════════════════════════╝")

	for {
		tampilkanMenu()
		pilihan := readInt("")

		switch pilihan {
		case 1:
			if n >= NMAX {
				fmt.Printf("\n✗ Data sudah penuh (maks %d).\n", NMAX)
			} else {
				inputAbsensi(&a, &n)
				fmt.Println("\n✓ Data berhasil diinput!")
			}

		case 2:
			if n == 0 {
				fmt.Println("\nBelum ada data. Silakan input terlebih dahulu.")
			} else {
				cetakAbsensi(a, n)
			}

		case 3:
			if n == 0 {
				fmt.Println("\nBelum ada data.")
			} else {
				cetakStatistik(a, n)
			}

		case 4:
			if n == 0 {
				fmt.Println("\nBelum ada data.")
			} else {
				cariEkstrimJam(a, n)
			}

		case 5:
			if n == 0 {
				fmt.Println("\nBelum ada data.")
			} else {
				keyword := readString("\nMasukkan nama yang dicari: ")
				sequentialSearchNama(a, n, keyword)
			}

		case 6:
			if n == 0 {
				fmt.Println("\nBelum ada data.")
			} else {
				prosedurBinarySearch(a, n)
			}

		case 7:
			if n == 0 {
				fmt.Println("\nBelum ada data.")
			} else {
				selectionSortAsc(&a, n)
				fmt.Println("\n✓ Data diurutkan berdasarkan ID (Ascending).")
				cetakAbsensi(a, n)
			}

		case 8:
			if n == 0 {
				fmt.Println("\nBelum ada data.")
			} else {
				selectionSortDesc(&a, n)
				fmt.Println("\n✓ Data diurutkan berdasarkan ID (Descending).")
				cetakAbsensi(a, n)
			}

		case 9:
			if n == 0 {
				fmt.Println("\nBelum ada data.")
			} else {
				insertionSortAsc(&a, n)
				fmt.Println("\n✓ Data diurutkan berdasarkan Jam Masuk (Ascending).")
				cetakAbsensi(a, n)
			}

		case 10:
			if n == 0 {
				fmt.Println("\nBelum ada data.")
			} else {
				insertionSortDesc(&a, n)
				fmt.Println("\n✓ Data diurutkan berdasarkan Jam Masuk (Descending).")
				cetakAbsensi(a, n)
			}

		case 0:
			fmt.Println("\nTerima kasih! Program selesai.")
			return

		default:
			fmt.Println("\n✗ Pilihan tidak valid. Masukkan angka 0-10.")
		}
	}
}