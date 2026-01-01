# go_photo

A simple CLI tool to organize photos from storage media into categorized folders.

The tool automatically creates `JPEG/`, `HEIF/`, and `RAW/` directories in your destination folder as needed.

## Usage

```bash
go_photo -src <source_dir> -dst <destination_dir>
```

### Examples

```bash
# Import from SD card
go_photo -src /Volumes/EOS_DIGITAL -dst ~/Pictures/Import

# Import from CFexpress card
go_photo -src /Volumes/CFEXPRESS -dst ~/Pictures/Import

# Import from external SSD
go_photo -src /Volumes/SanDisk -dst ~/Pictures/Import
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

Extensions are lowercased internally to avoid case sensitivity issues.

### Example: Source to Destination

**src /Volumes/XS20/DCIM**

The program walks all sub directories for you.

**Source (SD card):**
```
/Volumes/XS20/DCIM/100FUJI/*.RAF
/Volumes/XS20/DCIM/101FUJI/*.RAF
/Volumes/XS20/DCIM/102FUJI/*.JPG
```

**dst ~/Documents/XS20**

The program auto imports file types to their respective sub directory.

If the directory already exists it will simply append new files.

If not, it will create the sub directories for you.

**Destination after running `go_photo -src /Volumes/XS20 -dst ~/Documents/XS20`:**
```
~/Documents/XS20/
├── JPEG/
│   └── *.JPG
└── RAW/
    └── *.RAF
```
