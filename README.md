# PassKit
This is a library for generating Apple Wallet PKPasses.

## How to use

This library was heavily inspired by [drallgood's jpasskit library](https://github.com/drallgood/jpasskit) which was
written in Java, so the objects and functions are very similar to the ones available on jpasskit.

### Define a pass

To define a pass you use the `Pass` struct, which represents
the [pass.json](https://developer.apple.com/documentation/walletpasses/pass) file. This struct is modeled as closely as
possible to the json file, so adding data is straightforward:

```go
c := passkit.NewBoardingPass(passkit.TransitTypeAir)

// Utility functions for adding fields to a pass
c.AddHeaderField(passkit.Field{
    Key: "your_head_key",
    Label: "your_displayable_head_label",
    Value:"value",
})
c.AddPrimaryFields(passkit.Field{
    Key: "your_prim_key",
    Label: "your_displayable_prim_label",
    Value:"value",
})
c.AddSecondaryFields(passkit.Field{
    Key: "your_sec_key",
    Label: "your_displayable_sec_label",
    Value:"value",
})
c.AddAuxiliaryFields(passkit.Field{
    Key: "your_aux_key",
    Label: "your_displayable_aux_label",
    Value:"value",
})
c.AddBackFields(passkit.Field{
    Key: "your_back_key",
    Label: "your_displayable_back_label",
    Value:"value",
})

boarding := time.Date(2023, 1, 1, 23, 50, 00, 00, time.UTC)

pass := passkit.Pass{
    FormatVersion:       1,
    TeamIdentifier:      "TEAMID",
    PassTypeIdentifier:  "pass.type.id",
    AuthenticationToken: "123141lkjdasj12314",
    OrganizationName:    "Your Organization",
    SerialNumber:        "1234",
    Description:         "test",
    BoardingPass:        c,
    Semantics: []passkit.SemanticTag{
        {
            AirlineCode:            "AA1234",
            TransitProvider:        "American Airlines",
            DepartureAirportCode:   "LAX",
            DepartureAirportName:   "Los Angeles International Airport",
            DepartureGate:          "28",
            DepartureTerminal:      "2",
            DestinationAirportCode: "LGA",
            DestinationAirportName: "LaGuardia Airport",
            DestinationGate:        "12",
            DestinationTerminal:    "1",
            TransitStatus:          "On Time",
            OriginalBoardingDate:   &boarding,
        },
    },
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

To create the pass structure you need a `PassTemplate` instance, either using streams (with `InMemoryPassTemplate`) or
files (with `FolderPassTemplate`).

#### Using files

To create the pass bundle create an instance of `FolderPassTemplate` using the absolute file path of the folder
containing the files:

```go
template := passkit.NewFolderPassTemplate("/home/user/pass")
```

All the files in the folder will be loaded exactly as provided.

#### Using streams (In Memory)

The second approach is more flexible, having the option of loading files using data streams or directly downloaded from
a public URL:

```go
template := passkit.NewInMemoryPassTemplate()

template.AddFileBytes(passkit.BundleThumbnail, bytes)
template.AddFileBytesLocalized(passkit.BundleIcon, "en", bytes)
err := template.AddFileFromURL(passkit.BundleLogo, "https://example.com/file.png")
err := template.AddFileFromURLLocalized(passkit.BundleLogo, "en", "https://example.com/file.png")
err := template.AddAllFiles("/home/user/pass")
```

**Note**: There are no checks that the contents of a provided file are valid. If a PDF file is provided, but is
referenced as icon.png, when viewing the pass on a device there will be issues. It also doesn't provide any
authentication for the downloads, so the resources used must be public for the download to work as expected. The 
downloads use a default `http.Client` without any SSL configuration, so if the download is from an HTTPS site the 
certificate must be valid or the download will fail. 

### Signing and zipping a pass

As all passes [need to be signed when bundled](https://developer.apple.com/documentation/walletpasses/building_a_pass)
we need to use a `Signer` instance. There are two types of signers:

* `FileBasedSigner`: uses a temp folder to store the signed zip file contents
* `MemoryBasedSigner`: keeps the signed zip file contents in memory

To use any of the `Signer` instances you need an instance of `SigningInformation` to load the certificates used to
generate the `signature`. There are two ways to obtain an instance. Either reading the certificates from the filesystem,
or from already loaded bytes in memory:

```go
// Using the certificate files from your filesystem
signInfo, err := passkit.LoadSigningInformationFromFiles("/home/user/pass_cert.p12", "pass_cert_password", "/home/user/AppleWWDRCA.cer")
// Using the bytes from the certificates, loaded from a database or vault, for example.
signInfo, err := passkit.LoadSigningInformationFromBytes(passCertBytes, "pass_cert_password", wwdrcaBytes)
```

**Note**: When loading the signing information errors will be returned if the certificates are invalid (expired, not
certificates, etc)

Finally, to create the signed pass bundle you use the `Pass`, `Signer`, `SigningInformation`, and `PassTemplate`
instances created previously, for example:

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

## Contributing

Right now I'm not really working on a project where this library is being actively used, so any bugs or weird behaviour
is hard for me to detect and fix. That's why this project is open to contributions, just make a Pull Request with fixes
or any other feature requests, and I will probably accept it and merge it (I will at least try to check the code).
