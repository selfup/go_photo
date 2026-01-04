# go_photo

A simple CLI tool to organize photos and videos from storage media into categorized folders.

The tool automatically creates `JPEG/`, `HEIF/`, `RAW/`, `MOV/`, `BRAW/`, and `MP4/` directories in your destination folder as needed.

## Usage

```bash
go_photo -src <source_dir> -dst <destination_dir>
```

### Examples

```bash
# Import from SD card
go_photo -src /Volumes/VolumeName -dst ~/Pictures/DestinationFolder

# Import from CFexpress card
go_photo -src /Volumes/VolumeName -dst ~/Pictures/DestinationFolder

# Import from external SSD
go_photo -src /Volumes/VolumeName -dst ~/Pictures/DestinationFolder
```

## Install (as a cli tool)

```bash
go install github.com/selfup/go_photo@latest
```

## Build (if cloned)

```bash
go build -o go_photo .
```

## Supported Sources

- SD Cards
- CFexpress Cards
- External SSDs
- Any directory on your machine

## Supported Formats

| Extension | Destination Folder |
|-----------|-------------------|
| .jpg, .jpeg | JPEG/ |
| .heif, .heic | HEIF/ |
| .raw, .arw, .raf, .nef | RAW/ |
| .mov | MOV/ |
| .braw | BRAW/ |
| .mp4 | MP4/ |

Extensions are lowercased internally to avoid case sensitivity issues.

### Example: Source to Destination

**src /Volumes/XS20/DCIM**

The program walks all sub directories for you.

**Source (SD card):**
```
/Volumes/XS20/DCIM/100FUJI/*.RAF
/Volumes/XS20/DCIM/101FUJI/*.RAF
/Volumes/XS20/DCIM/102FUJI/*.JPG
/Volumes/BMPCC4K_SSD/*.BRAW
/Volumes/FX3/PRIVATE/M4ROOT/CLIP/*.MP4
```

**dst ~/Documents/XS20**

The program auto imports file types to their respective sub directory.

If the directory already exists it will simply append new files.

If not, it will create the sub directories for you.

**Destination after running `go_photo -src /Volumes/MediaName -dst ~/Documents/DestinationFolder`:**
```
~/Documents/DestinationFolder/
├── JPEG/
│   └── *.JPG (Ex: from XS20)
├── RAW/
│   └── *.RAF (Ex: from XS20)
├── MOV/
│   └── *.MOV (Ex: from XS20)
├── BRAW/
│   └── *.BRAW (Ex: from BMPCC4K)
└── MP4/
    └── *.MP4 (Ex: from FX3)
```

Real Examples:

`go run main.go -src E:/XS20/DCIM -dst D:/XS20`

`go run main.go -src /Volumes/OWC1TB/DCIM -dst /Volumes/T9/XH2S`

`go run main.go -src /Volumes/SanDisk2TB -dst /Volumes/T9/BMPCC4k`

`go run main.go -src /Volumes/Lexar512GB -dst /Volumes/T9/Pyxis6k`

`go run main.go -src /Volumes/FX3/PRIVATE/M4ROOT/CLIP -dst /Volumes/T9/FX3`
