import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";

const siteBucket = new aws.s3.Bucket("siteBucket", {});
const testFileAsset = new aws.s3.BucketObject("testFileAsset", {
    bucket: siteBucket,
    source: new pulumi.asset.FileAsset("file.txt"),
});
const testStringAsset = new aws.s3.BucketObject("testStringAsset", {
    bucket: siteBucket,
    source: new pulumi.asset.StringAsset("<h1>File contents</h1>"),
});
const testRemoteAsset = new aws.s3.BucketObject("testRemoteAsset", {
    bucket: siteBucket,
    source: new pulumi.asset.RemoteAsset("https://pulumi.test"),
});
const testFileArchive = new aws.s3.BucketObject("testFileArchive", {
    bucket: siteBucket,
    source: new pulumi.asset.FileArchive("file.tar.gz"),
});
const testRemoteArchive = new aws.s3.BucketObject("testRemoteArchive", {
    bucket: siteBucket,
    source: new pulumi.asset.RemoteArchive("https://pulumi.test/foo.tar.gz"),
});
const testAssetArchive = new aws.s3.BucketObject("testAssetArchive", {
    bucket: siteBucket,
    source: new pulumi.asset.AssetArchive({
        "file.txt": new pulumi.asset.FileAsset("file.txt"),
        "string.txt": new pulumi.asset.StringAsset("<h1>File contents</h1>"),
        "remote.txt": new pulumi.asset.RemoteAsset("https://pulumi.test"),
        "file.tar": new pulumi.asset.FileArchive("file.tar.gz"),
        "remote.tar": new pulumi.asset.RemoteArchive("https://pulumi.test/foo.tar.gz"),
        ".nestedDir": new pulumi.asset.AssetArchive({
            "file.txt": new pulumi.asset.FileAsset("file.txt"),
            "string.txt": new pulumi.asset.StringAsset("<h1>File contents</h1>"),
            "remote.txt": new pulumi.asset.RemoteAsset("https://pulumi.test"),
            "file.tar": new pulumi.asset.FileArchive("file.tar.gz"),
            "remote.tar": new pulumi.asset.RemoteArchive("https://pulumi.test/foo.tar.gz"),
        }),
    }),
});