class Fu < Formula
  desc "intuitive alternative to the find utility"
  homepage "https://github.com/kbrgl/fu"
  url "https://github.com/kbrgl/fu/archive/v2.0.0.tar.gz"
  sha256 "c9b3de10807f10da44ad5dda0a421d5b8ad1005964ba7d6ff8df7a270511de53"
  head "https://github.com/kbrgl/fu.git"

  depends_on "go" => :build

  def self.dependencies
    %w[
      github.com/alecthomas/kingpin
      github.com/kbrgl/fu/matchers
      github.com/kbrgl/fu/shallowradix
      github.com/mattn/go-isatty
      github.com/stretchr/powerwalk
    ]
  end

  def install
    # Ensure packages are installed to the right location
    ENV["GOPATH"] = buildpath
    puts ENV["GOPATH"]

    # Install dependencies
    self.class.dependencies.each do |dep|
      system "go", "get", dep
    end

    system "go", "build", "-o", "fu"
    bin.install "fu"
  end

  test do
    system "#{bin}/fu", "--version"
  end
end
