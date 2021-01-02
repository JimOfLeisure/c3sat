$Signature = @'
[DllImport(@"./decompressciv3.dylib")]
public static extern void HelloDll();
[DllImport(@"./decompressciv3.dylib")]
public static extern char* ReadFile(char* path);
[DllImport(@"./decompressciv3.dylib")]
public static extern void Decompress();
'@

$Type = Add-Type -MemberDefinition $Signature -Name Win32Utils -Namespace GoDllTest -PassThru

$Type::HelloDll()
$Type::ReadFile("foo")
#$Type::Decompress()
