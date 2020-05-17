# Quick Script to modify Apolyton's BIC format docs to Go struct member names

param (
    [Parameter(Mandatory=$true)]$String
)

$String.Split("`n") |
    %{ $PSItem.trim() } |
    Select-String -NotMatch "^$" | 
    %{ $PSItem.Line.Split("`t")[-1] } |
    %{ (Get-Culture).TextInfo.ToTitleCase($PSItem) -replace " |\(|\)|:" } |
    %{ $PSItem + " int32" } |
    Set-Clipboard
