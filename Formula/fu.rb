class Fu < Formula
  desc "intuitive alternative to the find utility"
  homepage "https://github.com/kbrgl/fu"
  url "https://github.com/kbrgl/fu/archive/v1.1.2.tar.gz"
  sha256 "34dd0eade8842c4b9e82b56e4b2154e9bcb34843d87aef568ea8fe16a1c1884a"
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
