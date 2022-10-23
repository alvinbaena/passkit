# PassKit
This is a library for generating Apple Wallet PKPasses.

## How to use
This library was heavily inspired by [drallgood's jpasskit library](https://github.com/drallgood/jpasskit) which was written in Java, so the objects and functions are very similar to the ones available on jpasskit.

### Define a pass
To define a pass you use the `Pass` struct, which represents the [pass.json](https://developer.apple.com/documentation/walletpasses/pass) file. This struct is modeled as closely as possible to the json file, so adding data is straightforward:

```go
c := passkit.NewBoardingPass(passkit.TransitTypeAir)
field := passkit.Field{
    Key: "key",
    Label: "label",
    Value:"value",
}

// Utility functions for adding fields to a pass
c.AddHeaderField(field)
c.AddPrimaryFields(field)
c.AddSecondaryFields(field)
c.AddAuxiliaryFields(field)
c.AddBackFields(field)

pass := passkit.Pass{
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
Usually, passes contain additional information that needs to be included in the final, signed pass, e.g.
* Images (icons, logos, background images)
* Translations

These templates are defined in the [apple wallet developer documentation](https://developer.apple.com/documentation/walletpasses/creating_the_source_for_a_pass).

To create the pass structure you need a `PassTemplate` instance, either using streams (with `InMemoryPassTemplate`) or files (with `FolderPassTemplate`).

#### Using files
To create the pass bundle create an instance of `FolderPassTemplate` using the absolute file path of the folder containing the files:

```go
template := passkit.NewFolderPassTemplate("/home/user/pass")
```

All the files in the folder will be loaded exactly as provided.

#### Using streams (In Memory)
The second approach is more flexible, having the option of loading files using data streams or directly downloaded from a public URL:

```go
template := passkit.NewInMemoryPassTemplate()

template.AddFileBytes(passkit.BundleThumbnail, bytes)
template.AddFileBytesLocalized(passkit.BundleIcon, "en", bytes)
err := template.AddFileFromURL(passkit.BundleLogo, "http://example.com/file.png")
err := template.AddFileFromURLLocalized(passkit.BundleLogo, "en", "http://example.com/file.png")
err := template.AddAllFiles("/home/user/pass")
```

**Note**: There are no checks that the contents of a provided file are valid. If a PDF file is provided, but is referenced as icon.png, when viewing the pass on a device there will be issues. It also doesn't provide any authentication for the downloads, so the resources used must be public for the download to work as expected.

### Signing and zipping a pass
As all passes [need to be signed when bundled](https://developer.apple.com/documentation/walletpasses/building_a_pass) we need to use a `Signer` instance. There are two types of signers:
* `FileBasedSigner`: uses a temp folder to store the signed zip file contents
* `MemoryBasedSigner`: keeps the signed zip file contents in memory

To use any of the `Signer` instances you need an instance of `SigningInformation` to load the certificates used to generate the `signature`. There are two ways to obtain an instance. Either reading the certificates from the filesystem, or from already loaded bytes in memory:
```go
signInfo, err := passkit.LoadSigningInformationFromFiles("/home/user/pass_cert.p12", "pass_cert_password", "/home/user/AppleWWDRCA.cer")
signInfo, err := passkit.LoadSigningInformationFromBytes(passCertBytes, "pass_cert_password", wwdrcaBytes)
```

**Note**: When loading the signing information errors will be returned if the certificates are invalid (expired, not certificates, etc)

Finally, to create the signed pass bundle you use the `Pass`, `Signer`, `SigningInformation`, and `PassTemplate` instances created previously, for example:

```go
signer := passkit.NewMemoryBasedSigner()
signInfo, err := passkit.LoadSigningInformationFromFiles("/home/user/pass_cert.p12", "pass_cert_password", "/home/user/AppleWWDRCA.cer")
if err != nil {
    panic(err)
}

z, err := signer.CreateSignedAndZippedPassArchive(&pass, template, signInfo)
if err != nil {
    panic(err)
}

err = os.WriteFile("/home/user/pass.pkpass", z, 0644)
if err != nil {
    panic(err)
}
```

After this step the pass bundle is ready to be distributed as you see fit.
