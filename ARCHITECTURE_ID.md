# Diagram Arsitektur protoc-gen-go-mapper

## Arsitektur

Diagram ini menunjukkan alur keseluruhan dari pemanggilan protoc hingga output kode yang dihasilkan. Plugin mengoordinasikan semua komponen, mengoordinasikan pemuatan konfigurasi, parsing proto, resolusi tipe, pendaftaran converter, pemuatan field handler, dan pembuatan kode.

**Komponen:**
- **protoc**: Kompilator Protocol Buffers yang memanggil plugin
- **Plugin**: Orkestrator utama yang mengoordinasikan semua komponen
- **Config**: Memuat dan memvalidasi konfigurasi mapper.yaml
- **Parser**: Mengekstrak file proto dan mengekstrak definisi message
- **Resolver**: Memetakan tipe protobuf ke tipe database spesifik
- **Registry**: Mengelola pendaftaran dan resolusi converter
- **HandlerRegistry**: Mengelola field handler untuk kasus khusus
- **Generator**: Menghasilkan kode pemetaan Go
- **Generated Code**: File output akhir dengan fungsi pemetaan

```mermaid
graph TD
    A[protoc] -->|Plugin Invocation| B[Plugin]
    B --> C[Config]
    B --> D[Parser]
    B --> E[Resolver]
    B --> F[Registry]
    B --> G[HandlerRegistry]
    B --> H[Generator]
    H --> I[Generated Code]
    
    C -->|mapper.yaml| J[Configuration]
    D -->|Proto File| K[Proto Descriptors]
    E -->|Type Mapping| L[Database Types]
    F -->|Converters| M[Converter Registry]
    G -->|Field Handlers| N[Handler Registry]
    
    style A fill:#e1f5ff
    style B fill:#fff4e1
    style I fill:#e8f5e9
```

## Diagram Interaksi Komponen

Diagram ini mengilustrasikan hubungan antara komponen inti dan pendukung. Plugin bertindak sebagai koordinator pusat, menghubungkan ke semua komponen inti (Parser, Resolver, Registry, Generator, HandlerRegistry) dan komponen Config. Komponen pendukung (Schema, Graph, Template, Converter, Handler) menyediakan definisi tipe, struktur pemetaan, template kode, logika konversi, dan kemampuan penanganan field.

**Komponen Inti:**
- **Plugin**: Koordinator pusat
- **Parser**: Mengekstrak struktur message dari file proto
- **Resolver**: Memetakan tipe antara protobuf dan database
- **Registry**: Mengelola pendaftaran converter
- **Generator**: Menghasilkan kode Go akhir
- **HandlerRegistry**: Mengelola logika kustom tingkat field

**Komponen Pendukung:**
- **Config**: Manajemen konfigurasi
- **Schema**: Definisi sistem tipe
- **Graph**: Struktur graf pemetaan
- **Template**: Template pembuatan kode
- **Converter**: Implementasi konversi tipe
- **Handler**: Implementasi field handler

```mermaid
graph LR
    subgraph "Core Components"
        A[Plugin]
        B[Parser]
        C[Resolver]
        D[Registry]
        E[Generator]
        F[HandlerRegistry]
    end
    
    subgraph "Supporting Components"
        G[Config]
        H[Schema]
        I[Graph]
        J[Template]
        K[Converter]
        L[Handler]
    end
    
    A --> B
    A --> C
    A --> D
    A --> E
    A --> F
    A --> G
    
    B --> H
    C --> H
    D --> K
    E --> J
    F --> L
    
    E --> I
    
    style A fill:#fff4e1
    style B fill:#e1f5ff
    style C fill:#e1f5ff
    style D fill:#e1f5ff
    style E fill:#e1f5ff
    style F fill:#e1f5ff
```

## Diagram Aliran Data

Diagram urutan ini menunjukkan alur eksekusi langkah demi langkah dari pemanggilan protoc hingga output file. Proses dimulai dengan protoc memanggil Plugin, yang kemudian memuat konfigurasi, mengurai file proto, menyelesaikan tipe, mendaftarkan converter, memuat field handler, membuat kode, dan akhirnya menulis kode yang dihasilkan ke file.

**Langkah Alur:**
1. **protoc** memanggil Plugin
2. **Plugin** memuat konfigurasi mapper.yaml
3. **Config** mengembalikan objek konfigurasi yang divalidasi
4. **Plugin** mengurai file proto untuk mengekstrak definisi message
5. **Parser** mengembalikan model skema dengan struktur message
6. **Plugin** menyelesaikan tipe protobuf ke tipe database
7. **Resolver** mengembalikan pemetaan tipe database
8. **Plugin** mendaftarkan converter bawaan dan kustom
9. **Registry** mengembalikan registry converter
10. **Plugin** memuat field handler dari konfigurasi
11. **HandlerRegistry** mengembalikan registry field handler
12. **Plugin** membuat kode pemetaan
13. **Generator** mengembalikan kode Go yang dihasilkan
14. **Plugin** menulis kode yang dihasilkan ke file

```mermaid
sequenceDiagram
    participant P as protoc
    participant PL as Plugin
    participant CFG as Config
    participant PR as Parser
    participant RES as Resolver
    participant REG as Registry
    participant HR as HandlerRegistry
    participant GEN as Generator
    participant OUT as Generated Code
    
    P->>PL: Invoke Plugin
    PL->>CFG: Load mapper.yaml
    CFG-->>PL: Config Object
    PL->>PR: Parse Proto File
    PR-->>PL: Schema Model
    PL->>RES: Resolve Types
    RES-->>PL: DB Type Mappings
    PL->>REG: Register Converters
    REG-->>PL: Converter Registry
    PL->>HR: Load Field Handlers
    HR-->>PL: Handler Registry
    PL->>GEN: Generate Code
    GEN-->>PL: Generated Go Code
    PL->>OUT: Write to File
```

## Arsitektur Resolver Database

Diagram ini menunjukkan bagaimana Resolver memetakan tipe protobuf ke tipe database spesifik berdasarkan backend database yang dikonfigurasi. Resolver menggunakan pernyataan switch untuk memilih resolver yang sesuai (SQLC, PGX, atau database_sql), yang masing-masing memiliki strategi pemetaan tipe sendiri.

**Backend Database:**
- **SQLC**: Menggunakan tipe pgtype untuk PostgreSQL (UUID, Timestamptz, Text, Numeric)
- **PGX**: Menggunakan tipe pgtype dengan variasi sedikit (UUID, Timestamp, Text, Numeric)
- **database_sql**: Menggunakan tipe Go standar (string, time.Time, sql.NullString, sql.NullTime)

**Pemetaan Tipe:**
- **UUID**: pgtype.UUID (SQLC/PGX) atau string (database_sql)
- **Timestamp**: pgtype.Timestamptz (SQLC), pgtype.Timestamp (PGX), atau time.Time (database_sql)
- **Text**: pgtype.Text (SQLC/PGX) atau string (database_sql)
- **Numeric**: pgtype.Numeric (SQLC/PGX) atau string (database_sql)

```mermaid
graph TD
    A[Resolver] --> B{Database Type}
    
    B -->|sqlc| C[SQLC Resolver]
    B -->|pgx| D[PGX Resolver]
    B -->|database_sql| E[Database SQL Resolver]
    
    C --> F[pgtype.UUID]
    C --> G[pgtype.Timestamptz]
    C --> H[pgtype.Text]
    C --> I[pgtype.Numeric]
    
    D --> F
    D --> J[pgtype.Timestamp]
    D --> H
    D --> I
    
    E --> K[string]
    E --> L[time.Time]
    E --> M[sql.NullString]
    E --> N[sql.NullTime]
    
    style A fill:#fff4e1
    style C fill:#e1f5ff
    style D fill:#e1f5ff
    style E fill:#e1f5ff
```

## Arsitektur Sistem Handler

Diagram ini menunjukkan sistem field handler yang menyediakan kustomisasi tingkat field yang fleksibel. HandlerRegistry mengelola beberapa tipe handler, masing-masing dengan tujuan spesifik. Handler diselesaikan menggunakan pencocokan berbasis prioritas, di mana handler dengan prioritas lebih tinggi diutamakan.

**Tipe Handler:**
- **type_assertion**: Menangani asersi tipe untuk field interface{} (misalnya, mengonversi interface{} ke []string untuk field array SQLC)
- **default_value**: Menetapkan nilai default untuk field yang tidak ada di sumber (misalnya, slice kosong untuk anak pohon)
- **skip**: Melewati field selama pemetaan (misalnya, field soft delete dalam respons)
- **field_mapping**: Menyediakan ekspresi kustom untuk kedua arah ToProto dan ToDB

**Level Prioritas:**
- **Prioritas 100**: Prioritas tertinggi (misalnya, handler skip)
- **Prioritas 75**: Prioritas menengah-tinggi (misalnya, handler nilai default)
- **Prioritas 50**: Prioritas menengah (misalnya, handler asersi tipe)

```mermaid
graph TD
    A[HandlerRegistry] --> B[Field Handlers]
    
    B --> C[type_assertion]
    B --> D[default_value]
    B --> E[skip]
    B --> F[field_mapping]
    
    C --> G[Interface to Type]
    D --> H[Set Default Value]
    E --> I[Skip Field]
    F --> J[Custom Expression]
    
    subgraph "Priority Resolution"
        K[Priority 100]
        L[Priority 75]
        M[Priority 50]
    end
    
    C --> M
    D --> L
    E --> K
    F --> M
    
    style A fill:#fff4e1
    style C fill:#e8f5e9
    style D fill:#e8f5e9
    style E fill:#fff3e0
    style F fill:#e8f5e9
```

## Sistem Konfigurasi

Diagram ini menunjukkan struktur file konfigurasi mapper.yaml. Objek Config memuat dan memvalidasi semua opsi konfigurasi, yang mengontrol bagaimana plugin menghasilkan fungsi pemetaan. Konversi tipe dan field handler mendukung fitur lanjutan seperti pencocokan pola, resolusi berbasis prioritas, dan strategi pointer.

**Opsi Konfigurasi:**
- **version**: Versi konfigurasi (harus "v1")
- **database**: Tipe database (sqlc, pgx, database_sql)
- **db_package**: Jalur paket Go untuk model database
- **package**: Nama paket Proto dan DB
- **type_mappings**: Pemetaan kustom message proto ke model DB
- **response_type_mappings**: Pemetaan message respons ke tipe Row SQLC
- **type_aliases**: Definisi konversi tipe yang dapat digunakan kembali
- **type_conversions**: Aturan konversi tipe kustom dengan pencocokan pola
- **field_handlers**: Logika kustom tingkat field
- **response_patterns**: Konfigurasi helper respons
- **pointer_settings**: Strategi penanganan pointer
- **messages**: Daftar message untuk membuat mapper

**Fitur Lanjutan:**
- **Pattern Matching**: Pencocokan nama field dan message berbasis regex
- **Priority**: Resolusi berbasis prioritas untuk beberapa kecocokan
- **Pointer Strategy**: strict, lenient, atau omit untuk field nullable

```mermaid
graph TD
    A[mapper.yaml] --> B[Config]
    
    B --> C[version]
    B --> D[database]
    B --> E[db_package]
    B --> F[package]
    B --> G[type_mappings]
    B --> H[response_type_mappings]
    B --> I[type_aliases]
    B --> J[type_conversions]
    B --> K[field_handlers]
    B --> L[response_patterns]
    B --> M[pointer_settings]
    B --> N[messages]
    
    J --> O[Pattern Matching]
    J --> P[Priority]
    J --> Q[Pointer Strategy]
    
    K --> R[Field Name]
    K --> S[DB Type]
    K --> T[Message]
    K --> U[Priority]
    
    style A fill:#e1f5ff
    style B fill:#fff4e1
    style J fill:#e8f5e9
    style K fill:#e8f5e9
```

## Registry Converter

Diagram ini menunjukkan sistem registry converter yang mengelola logika konversi tipe. Registry mempertahankan koleksi converter untuk pasangan tipe yang berbeda (Scalar, UUID, Timestamp, Decimal, Enum, Nullable, Slice, Message). Ketika type_conversions kosong (mode zero-config), converter generik (ConvertUUID, ConvertTimestamp, ConvertText) secara otomatis digunakan. Registry menggunakan resolusi berbasis prioritas untuk memilih converter terbaik ketika beberapa converter cocok dengan pasangan tipe.

**Converter Bawaan:**
- **ScalarConverter**: Menangani tipe skalar dasar (int32, int64, string, bool, float64)
- **UUIDConverter**: Menangani konversi tipe UUID
- **TimestampConverter**: Menangani konversi tipe timestamp
- **DecimalConverter**: Menangani konversi tipe desimal/numerik
- **EnumConverter**: Menangani konversi tipe enum
- **NullableConverter**: Menangani konversi tipe nullable/opsional
- **SliceConverter**: Menangani konversi tipe array/slice
- **MessageConverter**: Menangani konversi tipe message bersarang

**Converter Generik (Mode Zero-Config):**
- **ConvertUUID**: Konversi otomatis UUID ↔ string
- **ConvertTimestamp**: Konversi otomatis Timestamp ↔ time.Time
- **ConvertText**: Konversi otomatis Text ↔ string

**Proses Resolusi:**
1. Cocokkan pasangan tipe dengan semua converter yang terdaftar
2. Dapatkan prioritas untuk setiap converter yang cocok
3. Pilih converter dengan prioritas tertinggi
4. Kembalikan error jika beberapa converter memiliki prioritas yang sama (pemetaan ambigu)

```mermaid
graph TD
    A[Registry] --> B[Converters]
    
    B --> C[ScalarConverter]
    B --> D[UUIDConverter]
    B --> E[TimestampConverter]
    B --> F[DecimalConverter]
    B --> G[EnumConverter]
    B --> H[NullableConverter]
    B --> I[SliceConverter]
    B --> J[MessageConverter]
    
    subgraph "Generic Converters"
        K[ConvertUUID]
        L[ConvertTimestamp]
        M[ConvertText]
    end
    
    B --> K
    B --> L
    B --> M
    
    subgraph "Priority Resolution"
        N[Match Type Pair]
        O[Get Priority]
        P[Select Best]
    end
    
    C --> N
    D --> N
    E --> N
    N --> O
    O --> P
    
    style A fill:#fff4e1
    style K fill:#e8f5e9
    style L fill:#e8f5e9
    style M fill:#e8f5e9
```

## Struktur Paket

Diagram ini menunjukkan organisasi paket keseluruhan dari proyek protoc-gen-go-mapper. Proyek dibagi menjadi empat direktori utama: cmd (antarmuka baris perintah), internal (implementasi inti), pkg (paket publik), dan examples (implementasi sampel).

**Struktur Direktori:**
- **cmd**: Berisi alat baris perintah protoc-gen-go-mapper utama
- **internal**: Paket implementasi inti (config, converter, generator, graph, handler, parser, registry, resolver, schema, template, plugin)
- **pkg**: Paket publik (converter, errors, naming, types)
- **examples**: Implementasi sampel (simple, medium, complex, advanced)

**Paket Internal:**
- **config**: Pemuatan dan validasi konfigurasi
- **converter**: Implementasi converter generik
- **generator**: Logika pembuatan kode
- **graph**: Struktur graf pemetaan
- **handler**: Implementasi field handler
- **parser**: Parsing file proto
- **registry**: Pendaftaran dan resolusi converter
- **resolver**: Resolusi tipe untuk database yang berbeda
- **schema**: Definisi sistem tipe
- **template**: Template pembuatan kode
- **plugin**: Orkestrasi plugin utama

**Paket Publik:**
- **converter**: Antarmuka converter publik
- **errors**: Definisi error
- **naming**: Utilitas konvensi penamaan
- **types**: Utilitas sistem tipe

```mermaid
graph TD
    A[protoc-gen-go-mapper] --> B[cmd]
    A --> C[internal]
    A --> D[pkg]
    A --> E[examples]
    
    B --> F[protoc-gen-go-mapper]
    
    C --> G[config]
    C --> H[converter]
    C --> I[generator]
    C --> J[graph]
    C --> K[handler]
    C --> L[parser]
    C --> M[registry]
    C --> N[resolver]
    C --> O[schema]
    C --> P[template]
    C --> Q[plugin]
    
    D --> R[converter]
    D --> S[errors]
    D --> T[naming]
    D --> U[types]
    
    E --> V[simple]
    E --> W[Medium]
    E --> X[Complex]
    E --> Y[Advanced]
    
    style A fill:#e1f5ff
    style C fill:#fff4e1
    style D fill:#e8f5e9
    style E fill:#f3e5f5
```

## Alur Mode Zero-Config

Diagram ini menunjukkan alur mode zero-config, yang diaktifkan ketika type_conversions kosong dalam konfigurasi. Dalam mode ini, plugin secara otomatis menggunakan converter generik bawaan untuk tipe umum (UUID, Timestamp, Text) tanpa memerlukan konfigurasi eksplisit. Converter generik diinlining langsung ke dalam kode yang dihasilkan, membuatnya mandiri tanpa dependensi eksternal.

**Proses Zero-Config:**
1. Periksa apakah type_conversions kosong
2. Jika kosong, aktifkan converter generik
3. Deteksi otomatis tipe field (UUID, Timestamp, Text)
4. Terapkan converter generik yang sesuai
5. Inline fungsi converter dalam kode yang dihasilkan
6. Buat kode pemetaan akhir

**Converter Generik:**
- **ConvertUUID**: Menangani konversi UUID ↔ string untuk field nullable dan non-nullable
- **ConvertTimestamp**: Menangani konversi Timestamp ↔ time.Time
- **ConvertText**: Menangani konversi Text ↔ string untuk field nullable dan non-nullable

**Manfaat:**
- Tidak ada konfigurasi yang diperlukan untuk tipe umum
- Kode yang dihasilkan mandiri
- Deteksi tipe otomatis
- Kompleksitas konfigurasi berkurang

```mermaid
graph TD
    A[Start] --> B{type_conversions empty?}
    B -->|Yes| C[Enable Generic Converters]
    B -->|No| D[Use Custom Conversions]
    
    C --> E[Auto-detect Types]
    E --> F[UUID Fields]
    E --> G[Timestamp Fields]
    E --> H[Text Fields]
    
    F --> I[ConvertUUID]
    G --> J[ConvertTimestamp]
    H --> K[ConvertText]
    
    I --> L[Inline in Generated Code]
    J --> L
    K --> L
    
    D --> M[Use Registry]
    M --> N[Apply Custom Logic]
    
    L --> O[Generate Code]
    N --> O
    
    style C fill:#e8f5e9
    style D fill:#fff3e0
    style L fill:#e1f5ff
```

## Pipeline Lengkap

Diagram ini menunjukkan pipeline end-to-end lengkap dari pemanggilan pengguna hingga output kode yang dihasilkan. Proses dimulai dengan pengguna menjalankan perintah protoc, yang memanggil Plugin. Plugin kemudian memuat dan memvalidasi konfigurasi, mengurai file proto, menyelesaikan tipe, mendaftarkan converter, memuat field handler, membuat fungsi pemetaan (ToProto, ToDB, dan Response Helpers), dan akhirnya menulis kode yang dihasilkan ke file .proto_mapper.pb.go.

**Tahap Pipeline:**
1. **Pemanggilan Pengguna**: Pengguna menjalankan protoc dengan plugin mapper
2. **Inisialisasi Plugin**: Plugin.New membuat instance plugin
3. **Pemuatan Konfigurasi**: Config.Load membaca mapper.yaml
4. **Validasi Konfigurasi**: Config.Validate memeriksa konfigurasi
5. **Generasi**: Plugin.Generate mengoordinasikan proses pembuatan
6. **Parsing Proto**: Parser.ParseFile mengekstrak definisi message
7. **Resolusi Tipe**: Resolver.Resolve memetakan tipe ke tipe database
8. **Pendaftaran Converter**: Registry.Register mendaftarkan converter
9. **Pemuatan Handler**: HandlerRegistry.Load memuat field handler
10. **Pembuatan Kode**: Generator.Generate menghasilkan fungsi pemetaan
11. **Output**: Kode yang dihasilkan ditulis ke file .proto_mapper.pb.go

**Fungsi yang Dihasilkan:**
- **Fungsi ToProto**: Mengonversi model database ke message protobuf
- **Fungsi ToDB**: Mengonversi message protobuf ke model database
- **Helper Respons**: Fungsi khusus untuk respons daftar

```mermaid
graph TD
    A[User] -->|protoc command| B[protoc]
    B -->|Plugin| C[Plugin.New]
    C -->|Load Config| D[Config.Load]
    D -->|Validate| E[Config.Validate]
    E -->|Valid| F[Plugin.Generate]
    
    F --> G[Parser.ParseFile]
    G --> H[Schema Model]
    
    F --> I[Resolver.Resolve]
    I --> J[DB Type Mappings]
    
    F --> K[Registry.Register]
    K --> L[Converter Registry]
    
    F --> M[HandlerRegistry.Load]
    M --> N[Field Handlers]
    
    F --> O[Generator.Generate]
    O --> P[ToProto Functions]
    O --> Q[ToDB Functions]
    O --> R[Response Helpers]
    
    P --> S[Generated Code]
    Q --> S
    R --> S
    
    S --> T[Write to File]
    T --> U[.proto_mapper.pb.go]
    
    style A fill:#e1f5ff
    style C fill:#fff4e1
    style O fill:#e8f5e9
    style U fill:#f3e5f5
```
