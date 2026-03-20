class Searchcli < Formula
  desc "Command-line web search tool with DuckDuckGo and Bing support"
  homepage "https://github.com/scian0204/SearchCli"
  url "https://github.com/scian0204/SearchCli/releases/download/v1.0.0/searchcli-darwin-arm64"
  version "1.0.0"
  sha256 "6633e15f996b8b8cc801af8bcc518cb616b97c86c0f9df2f58cd83b2362cd5d2"
  license "MIT"

  depends_on :macos

  def install
    bin.install "searchcli-darwin-arm64" => "searchcli"
  end

  test do
    system "#{bin}/searchcli", "-help"
  end
end
