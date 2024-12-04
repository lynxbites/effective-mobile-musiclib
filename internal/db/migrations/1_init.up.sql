CREATE TABLE songs (
    songId          SERIAL PRIMARY KEY,
    groupName       TEXT,
    songName        TEXT,
    releaseDate     DATE,
    songLyrics      TEXT,
    songLink        TEXT
);

INSERT INTO songs (groupName, songName, releaseDate, songLyrics, songLink) VALUES ('B group', 'D name', '2000-01-01', 'text', 'link');
INSERT INTO songs (groupName, songName, releaseDate, songLyrics, songLink) VALUES ('A group', 'F name', '2000-01-02', 'text', 'link');
INSERT INTO songs (groupName, songName, releaseDate, songLyrics, songLink) VALUES ('F group', 'E name', '2000-01-05', 'text', 'link');
INSERT INTO songs (groupName, songName, releaseDate, songLyrics, songLink) VALUES ('D group', 'A name', '2000-01-04', 'text', 'link');
INSERT INTO songs (groupName, songName, releaseDate, songLyrics, songLink) VALUES ('C group', 'B name', '2000-01-03', 'text', 'link');
INSERT INTO songs (groupName, songName, releaseDate, songLyrics, songLink) VALUES ('E group', 'C name', '2000-01-06', 'text', 'link');