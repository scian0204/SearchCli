class Searchcli < Formula
  desc "Command-line web search tool with DuckDuckGo and Bing support"
  homepage "https://github.com/scian0204/SearchCli"
  url "https://github.com/scian0204/SearchCli/releases/download/v1.0.0/searchcli"
  version "1.0.0"
  sha256 ":placeholder"
  license "MIT"

  depends_on :macos

  def install
    bin.install "searchcli" => "searchcli"
  end

  test do
    system "#{bin}/searchcli", "-help"
  end
end
