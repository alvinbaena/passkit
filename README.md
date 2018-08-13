# PassKit
This is a library for generating Apple Wallet PKPasses.

## How to use
This library was heavily inspired by [drallgood's jpasskit library](https://github.com/drallgood/jpasskit) which was written in Java, so the objects and functions are very similar to the ones available on jpasskit.

### Define a pass
To define a pass you use the `Pass` struct, which represents the [pass.json](https://developer.apple.com/library/archive/documentation/UserExperience/Reference/PassKit_Bundle/Chapters/TopLevel.html#//apple_ref/doc/uid/TP40012026-CH2-SW1) file. This struct is modeled as closely as possible to the json file, so adding data is straightforward:

```go
c := passkit.NewBoardingPass(passkit.TransitTypeAir)
field := passkit.Field{
    Key: "key",
    Label: "label",
    Value:"value",
}

c.AddHeaderField(field)
c.AddPrimaryFields(field)
c.AddSecondaryFields(field)

p := passkit.Pass{
    FormatVersion:       1,
    TeamIdentifier:      "TEAMID",
    PassTypeIdentifier:  "pass.type.id",
    AuthenticationToken: "123141lkjdasj12314",
    OrganizationName:    "Your Organization",
    SerialNumber:        "1234",
    Description:         "test",
    BoardingPass:         c,
    Barcodes: []passkit.Barcode{
        {
            Format:          passkit.BarcodeFormatPDF417,
            Message:         "1312312312312312312312312312",
            MessageEncoding: "utf-8",
        },
    },
}
```

### Templates
Usually, passes contain additional information that need to be included in the final, signed pass, e.g.
* Images (icons, logos, background images)
* Translations

These templates are defined in the [apple wallet developer documentation](https://developer.apple.com/library/archive/documentation/UserExperience/Reference/PassKit_Bundle/Chapters/PackageStructure.html#//apple_ref/doc/uid/TP40012026-CH1-SW1).

To create the pass structure you need a `PassTemplate` instance, either using streams (with `InMemoryPassTemplate`) or files (with `FolderPassTemplate`).

#### Using files
To load the pass with files in the file system you create an instance of `FolderPassTemplate` passing the absolute file path of the folder:

```go
folderTemplate := passkit.NewFolderPassTemplate("/home/user/pass")
```

When building the pass the files in the folder will be added.

#### Using streams (In Memory)
The second approach is more flexible, having the option of loading files using data streams:

```go
memTemplate := passkit.NewInMemoryPassTemplate()

memTemplate.AddFileBytes(passkit.BundleThumbnail, bytes)
memTemplate.AddFileBytesLocalized(passkit.BundleIcon, "en", bytes)
err := memTemplate.AddFileFromURL(passkit.BundleLogo, "http://example.com/file.png")
err := memTemplate.AddFileFromURLLocalized(passkit.BundleLogo, "en", "http://example.com/file.png")
err := memTemplate.AddAllFiles("/home/user/pass")
```
**Note**: There are no checks, that the content of a provided file is valid. So if you'd provide a PDF file but store it as icon.png, it will not work.

### Signing and zipping a pass
To create a pkpass file you need to use a `Signer`. There are two types of signers:
* FileBasedSigner (uses a temp folder to create the zip file)
* MemoryBasedSigner (creates the zip on memory as bytes)

To use any of the `Signer` instances you need an instance of `SigningInformation`. There are two methods to obtain an instance:
```go
signer, err := passkit.LoadSigningInformationFromFiles("/home/user/pass_cert.p12", "password", "/home/user/AppleWWDRCA.cer")
signer, err := passkit.LoadSigningInformationFromBytes(passCertBytes, "password", wwdrcaBytes)
```

To create a zip you use the `Pass`, `Signer`, `SigningInformation`, and `PassTemplate` instances created previously:

```go
signer := passkit.NewMemoryBasedSigner()
z, err := signer.CreateSignedAndZippedPassArchive(&pass, template, signInfo)
if err != nil {
    panic(err)
}

err = ioutil.WriteFile("/home/user/pass.pkpass", z, 0644)
if err != nil {
    panic(err)
}
```