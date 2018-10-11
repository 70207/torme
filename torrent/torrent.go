package torrent



type TorrentInfo struct {
    Encoding         string
    Name             string
    PieceLength      int
    Pieces       [][]byte
    Files          []FileInfo
    Publisher        string
    PublisherUrl     string
    Announces      []AnnounceList
    CreatedBy        string
    CreationDate     int
    Comment          string
}

type FileInfo struct{
    Path           []string
    Length           int64
}

type AnnounceList struct{
    Announce       []string
}