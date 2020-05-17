 ri .\civ3sat.exe
 go build .
 if ($LASTEXITCODE) {break}
 #$Result = 
 gci C:\temp\saves\*.sav |
    %{
        $filename = $psitem.Name
        .\civ3sat.exe z $PSItem |
            ConvertFrom-Json #|
            # Add-Member -Type NoteProperty -Name "FileName" -Value $filename -PassThru
    }
    #$Result

break

ri .\civ3sat.exe; go build .; if ($LASTEXITCODE) {break}; gci C:\temp\saves\*.sav | %{ $filename = $psitem.Name; .\civ3sat.exe z $PSItem | ConvertFrom-Json}
