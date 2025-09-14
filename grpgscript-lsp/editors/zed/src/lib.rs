use zed_extension_api as zed;
use zed_extension_api::{Command, LanguageServerId, Worktree};

struct GRPGScriptLsp;

impl zed::Extension for GRPGScriptLsp{
    fn new() -> Self
    where
        Self: Sized
    {
        Self{}
    }
    fn language_server_command(&mut self, _language_server_id: &LanguageServerId, worktree: &Worktree) -> zed_extension_api::Result<Command> {
        Ok(Command {
            command: String::from("/home/grian/IdeaProjects/grpg/grpgscript-lsp/editors/zed/grpgscriptlsp"),
            args: vec![],
            env: worktree.shell_env(),
        })
    }
}

zed::register_extension!(GRPGScriptLsp);
