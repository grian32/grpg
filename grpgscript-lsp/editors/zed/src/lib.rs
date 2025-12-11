use zed_extension_api as zed;
use zed_extension_api::{Command, LanguageServerId, Worktree};

struct GRPGScriptLsp {}

impl zed::Extension for GRPGScriptLsp{
    fn new() -> Self
    where
        Self: Sized
    {
        zed::download_file(
            "http://51.83.129.212:4022/assets/grpgscriptlsp",
            "./grpgscriptlsp",
            zed_extension_api::DownloadedFileType::Uncompressed,
        ).expect("grpgscriptlsp: failed to download lsp binary");

        zed::make_file_executable("./grpgscriptlsp").expect("grpgscriptlsp: failed to make lsp binary executable");

        Self {  }
    }
    fn language_server_command(&mut self, _language_server_id: &LanguageServerId, worktree: &Worktree) -> zed_extension_api::Result<Command> {
        Ok(Command {
            command: String::from("./grpgscriptlsp"),
            args: vec![],
            env: worktree.shell_env(),
        })
    }
}

zed::register_extension!(GRPGScriptLsp);
