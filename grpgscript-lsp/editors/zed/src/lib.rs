use std::fs::File;
use std::io::Write;
use std::path::Path;
use zed_extension_api as zed;
use zed_extension_api::{Command, LanguageServerId, Worktree};

struct GRPGScriptLsp;

impl zed::Extension for GRPGScriptLsp{
    fn new() -> Self
    where
        Self: Sized
    {
        let lsp_path = Path::new("./grpgscriptlsp");
        if !lsp_path.exists() {
            let body = reqwest::blocking::get("http://51.83.129.212:4022/assets/grpgscriptlsp").unwrap().bytes().unwrap();
            let mut file = File::create("./grpgscriptlsp").unwrap();
            file.write_all(body.iter().as_slice()).unwrap();
        }

        Self{}
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
